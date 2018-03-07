package batexpe

import (
	"bytes"
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(testTimeout)*time.Second)
	defer cancel()

	termination := make(chan RobinResult)
	go executeRobinInnerCtx(ctx, descriptionFile, termination)

	rresult := <-termination
	return rresult
}

func executeRobinInnerCtx(ctx context.Context, descriptionFile string,
	onexit chan RobinResult) {
	cmd := exec.Command("robin")
	cmd.Args = []string{"robin", "--json-logs", descriptionFile}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	log.WithFields(log.Fields{
		"command": cmd,
		"context": ctx,
	}).Debug("Starting robin")

	if err := cmd.Start(); err != nil {
		log.WithFields(log.Fields{
			"command": cmd,
		}).Fatal("Could not start robin")
	}

	var rresult RobinResult
	rresult.Finished = false
	rresult.ReturnCode = -1

	if err := cmd.Wait(); err != nil {

		if ctx.Err() != nil {
			log.WithFields(log.Fields{
				"err":     ctx.Err(),
				"command": cmd.Args,
			}).Info("Test timeout reached!")
		} else {
			rresult.Finished = true
		}

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
		rresult.Finished = true
		rresult.ReturnCode = 0
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

func WasBatsimSuccessful(robinJsonLines []interface{}) bool {
	for _, object := range robinJsonLines {
		line_as_map := object.(map[string]interface{})

		if line_as_map["msg"] == "Simulation subprocess succeeded" &&
			line_as_map["process name"] == "Batsim" {
			return true
		} else if line_as_map["msg"] == "Simulation subprocess failed" &&
			line_as_map["process name"] == "Batsim" {
			return false
		}
	}

	return false
}

func WasSchedSuccessful(robinJsonLines []interface{}) (successful, present bool) {
	present = false
	for _, object := range robinJsonLines {
		line_as_map := object.(map[string]interface{})

		if line_as_map["msg"] == "Starting simulation" {
			_, sched_in_simu := line_as_map["scheduler command"]
			present = sched_in_simu
		} else if line_as_map["msg"] == "Simulation subprocess succeeded" &&
			line_as_map["process name"] == "Scheduler" {
			return true, present
		} else if line_as_map["msg"] == "Simulation subprocess failed" &&
			line_as_map["process name"] == "Scheduler" {
			return false, present
		}
	}

	return false, present
}

func WasContextClean(robinJsonLines []interface{}) bool {
	for _, object := range robinJsonLines {
		line_as_map := object.(map[string]interface{})

		if line_as_map["msg"] == "Starting simulation" {
			return true
		} else if line_as_map["msg"] == "Context remains invalid" {
			return false
		}
	}

	return false
}
