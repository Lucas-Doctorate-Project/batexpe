package batexpe

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type RobinResult struct {
	Finished   bool
	ReturnCode int
	Output     string
}

func RunRobin(descriptionFile string, testTimeout float64) RobinResult {
	termination := make(chan RobinResult)
	go executeRobinWithTimeout(testTimeout, descriptionFile, termination)

	rresult := <-termination
	return rresult
}

func executeRobinWithTimeout(timeout float64, descriptionFile string,
	onexit chan RobinResult) {
	cmd := exec.Command("robin")
	cmd.Args = []string{"robin", "--json-logs", descriptionFile}

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
		rresult.ReturnCode = -1
		KillProcess(robinPid)
	case err := <-done:
		rresult.Finished = true
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				// Exited with non-zero exit code
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					rresult.ReturnCode = status.ExitStatus()
				} else {
					log.WithFields(log.Fields{
						"command": cmd,
						"err":     err,
					}).Fatal("Cannot retrieve robin exit code (case 1)")
				}
			} else {
				log.WithFields(log.Fields{
					"command": cmd,
					"err":     err,
				}).Fatal("Cannot retrieve robin exit code (case 2)")
			}
		} else {
			rresult.ReturnCode = 0
		}
	}

	rresult.Output = stdout.String()
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
		if err := json.Unmarshal([]byte(lines[i]), &jsonLines[i]); err != nil {
			return nil, fmt.Errorf("Could not unmarshall JSON line: %s", lines[i])
		}
	}

	return jsonLines, nil
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
