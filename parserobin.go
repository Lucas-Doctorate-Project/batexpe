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
				"command": cmd,
				"context": ctx,
			}).Error("Test timeout reached!")
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

	json_lines := make([]interface{}, len(lines))

	for i := 0; i < len(lines); i++ {
		log.WithFields(log.Fields{
			"line": lines[i],
		}).Debug("Parsing line")
		if err := json.Unmarshal([]byte(lines[i]), &json_lines[i]); err != nil {
			return nil, fmt.Errorf("Could not unmarshall JSON line: %s", lines[i])
		}
	}

	return json_lines, nil
}
