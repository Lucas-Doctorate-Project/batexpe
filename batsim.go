package batexpe

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

type BatsimArgs struct {
	Socket       string
	ExportPrefix string
	BatexecMode  bool
}

func ParseBatsimCommand(batcmd string) (batargs BatsimArgs, err error) {
	// Goal: Run batsim with --dump-execution-context and parse its JSON output.
	executedCmd := batcmd + " --dump-execution-context"

	// Create a temporary command file for the command to run.
	cmdFile, err := ioutil.TempFile(".", "tmp-dump-batcmd")
	if err != nil {
		return batargs, fmt.Errorf("Cannot create temporary file")
	}
	defer os.Remove(cmdFile.Name())

	// Dump command to the temporary file.
	cmdFile.Write([]byte(executedCmd))

	cmd := exec.Command("bash")
	cmd.Args = []string{cmd.Args[0], "-eux", cmdFile.Name()}

	log.WithFields(log.Fields{
		"command": executedCmd,
	}).Debug("Executing Batsim to parse its arguments")

	out, err := cmd.Output()
	if err != nil {
		return batargs, err
	}

	log.WithFields(log.Fields{
		"output": string(out),
	}).Debug("Batsim execution was successful (to parse its arguments)")

	// Parse output JSON to know the command execution context
	var jsonData map[string]interface{}
	if err := json.Unmarshal(out, &jsonData); err != nil {
		return batargs, fmt.Errorf("Could not parse Batsim (expected to be JSON) output")
	}

	batargs.Socket = jsonData["socket_endpoint"].(string)
	batargs.ExportPrefix = jsonData["export_prefix"].(string)
	batargs.BatexecMode = !jsonData["external_scheduler"].(bool)

	return batargs, nil
}
