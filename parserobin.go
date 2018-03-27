package batexpe

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RobinResult struct {
	Finished  bool
	Succeeded bool
	Output    string
}

func RunRobin(descriptionFile, coverFile string,
	testTimeout float64) RobinResult {
	termination := make(chan RobinResult)
	go executeRobinWithTimeout(testTimeout, descriptionFile, coverFile,
		termination)

	rresult := <-termination
	return rresult
}

func executeRobinWithTimeout(timeout float64, descriptionFile,
	coverFile string, onexit chan RobinResult) {
	cmd := exec.Command("robin")

	if coverFile == "" {
		cmd.Args = []string{"robin", "--json-logs", descriptionFile}
	} else {
		testArg := "-test.coverprofile=" + coverFile
		cmd = exec.Command("robin.cover")
		cmd.Args = []string{"robin.cover", testArg, descriptionFile,
			"--json-logs"}
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	log.WithFields(log.Fields{
		"command": cmd,
		"timeout": timeout,
	}).Debug("Starting robin")

	if err := cmd.Start(); err != nil {
		log.WithFields(log.Fields{
			"command": cmd,
		}).Fatal("Could not start robin")
	}

	robinPid := cmd.Process.Pid
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	var rresult RobinResult

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		log.WithFields(log.Fields{
			"command": cmd.Args,
			"timeout": timeout,
		}).Info("Test timeout reached!")
		rresult.Finished = false
		rresult.Succeeded = false

		KillProcess(robinPid)
		<-done
		rresult.Output = stdout.String()
	case <-done:
		rresult.Output = stdout.String()
		rresult.Finished = true

		if coverFile == "" {
			// robin is directly executed, its return code can be retrieved
			rresult.Succeeded = cmd.ProcessState.Success()
		} else {
			// robin.cover is called. It should always return 0.
			// Robin's return code should be written in the program output
			returnCode, err :=
				retrieveRobinReturnCodeInRobincoverOutput(rresult.Output)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Fatal("Cannot retrieve return code from robin.cover output")
			}
			rresult.Succeeded = returnCode == 0
		}
	}
	onexit <- rresult
}

func ParseRobinOutput(output string) ([]interface{}, error) {
	splitFn := func(c rune) bool {
		return c == '\n'
	}
	lines := strings.FieldsFunc(output, splitFn)

	jsonLines := make([]interface{}, len(lines))

	for i := 0; i < len(lines); i++ {
		log.WithFields(log.Fields{
			"line": lines[i],
		}).Debug("Parsing line")

		// Parse line if it is not a coverage print
		if lines[i] != "PASS" && !strings.HasPrefix(lines[i], "cover") &&
			!strings.HasPrefix(lines[i], "Robin return code:") {
			if err := json.Unmarshal([]byte(lines[i]), &jsonLines[i]); err != nil {
				return nil, fmt.Errorf("Could not unmarshall JSON line: %s", lines[i])
			}
		}
	}

	return jsonLines, nil
}

func retrieveRobinReturnCodeInRobincoverOutput(output string) (int, error) {
	r := regexp.MustCompile(`Robin return code:\s*(?P<returnCode>\d+)\s*`)

	match := r.FindStringSubmatch(output)
	if match == nil {
		return -1, fmt.Errorf("Return line not found")
	}

	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	returnCode, err := strconv.ParseInt(result["returnCode"], 10, 32)
	if err != nil {
		return -1, fmt.Errorf("Cannot convert return code %v to int",
			result["returnCode"])
	}
	return int(returnCode), nil
}

func WasBatsimSuccessful(robinJsonLines []interface{}) (successful, killed bool) {
	for _, object := range robinJsonLines {
		lineAsMap := object.(map[string]interface{})

		if lineAsMap["msg"] == "Simulation subprocess succeeded" &&
			lineAsMap["process name"] == "Batsim" {
			return true, false
		} else if lineAsMap["msg"] == "Simulation subprocess failed" &&
			lineAsMap["process name"] == "Batsim" {
			batsimKilled := strings.HasPrefix(lineAsMap["err"].(string),
				"signal: ")
			return false, batsimKilled
		} else if lineAsMap["msg"] == "Simulation subprocess failed (simulation timeout reached)" &&
			lineAsMap["process name"] == "Batsim" {
			return false, true
		}
	}

	return false, false
}

func WasSchedSuccessful(robinJsonLines []interface{}) (successful, present, killed bool) {
	present = false
	for _, object := range robinJsonLines {
		lineAsMap := object.(map[string]interface{})

		if lineAsMap["msg"] == "Starting simulation" {
			_, sched_in_simu := lineAsMap["scheduler command"]
			present = sched_in_simu
		} else if lineAsMap["msg"] == "Simulation subprocess succeeded" &&
			lineAsMap["process name"] == "Scheduler" {
			return true, present, false
		} else if lineAsMap["msg"] == "Simulation subprocess failed" &&
			lineAsMap["process name"] == "Scheduler" {
			schedKilled := strings.HasPrefix(lineAsMap["err"].(string),
				"signal: ")
			return false, present, schedKilled
		} else if lineAsMap["msg"] == "Simulation subprocess failed (simulation timeout reached)" &&
			lineAsMap["process name"] == "Scheduler" {
			return false, present, true
		}
	}

	return false, present, false
}

func WasContextClean(robinJsonLines []interface{}) bool {
	for _, object := range robinJsonLines {
		lineAsMap := object.(map[string]interface{})

		if lineAsMap["msg"] == "Starting simulation" {
			return true
		} else if lineAsMap["msg"] == "Context remains invalid" {
			return false
		}
	}

	return false
}
