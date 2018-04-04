package batexpe

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	SUCCESS int = iota
	TIMEOUT
	FAILURE
)

type CmdFinishedMsg struct {
	Name  string
	State int
}

func PrepareDirs(exp Experiment) {
	// Create output directory if needed
	err := CreateDirIfNeeded(exp.OutputDir)
	if err != nil {
		log.WithFields(log.Fields{
			"directory": exp.OutputDir,
			"err":       err,
		}).Fatal("Cannot create output directory")
	}

	err = CreateDirIfNeeded(exp.OutputDir + "/log")
	if err != nil {
		log.WithFields(log.Fields{
			"directory": exp.OutputDir + "/log",
			"err":       err,
		}).Fatal("Cannot create log directory")
	}

	err = CreateDirIfNeeded(exp.OutputDir + "/cmd")
	if err != nil {
		log.WithFields(log.Fields{
			"directory": exp.OutputDir + "/cmd",
			"err":       err,
		}).Fatal("Cannot create command directory")
	}
}

func waitReadyForSimulation(exp Experiment, batargs BatsimArgs) {
	port, err := PortFromBatSock(batargs.Socket)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"extracted socket endpoint": batargs.Socket,
			"batsim command":            exp.Batcmd,
		}).Fatal("Cannot retrieve port from Batsim socket")
	}

	log.WithFields(log.Fields{
		"ready timeout (seconds)":   exp.ReadyTimeout,
		"extracted port":            port,
		"extracted socket endpoint": batargs.Socket,
		"batsim command":            exp.Batcmd,
	}).Info("Waiting for valid context")

	socketInUse := true
	anotherBatsim := true

	sockChan := make(chan int)
	batChan := make(chan int)

	go waitTcpPortAvailableSs(port, sockChan)
	go waitNoConflictingBatsim(port, batChan)

	for socketInUse || anotherBatsim {
		select {
		case <-time.After(time.Duration(exp.ReadyTimeout) * time.Second):
			log.WithFields(log.Fields{
				"ready timeout (seconds)":    exp.ReadyTimeout,
				"scanned port":               port,
				"batsim command":             exp.Batcmd,
				"socket in use":              socketInUse,
				"conflicting batsim running": anotherBatsim,
			}).Fatal("Context remains invalid")
		case <-sockChan:
			socketInUse = false
		case <-batChan:
			anotherBatsim = false
		}
	}
}

func waitNoConflictingBatsim(port uint16, onexit chan int) {
	r := regexp.MustCompile(`(?m)^\S*\bbatsim\s+.*$`)
	for {
		// Retrieve running Batsim processes
		psCmd := exec.Command("ps")
		psCmd.Args = []string{psCmd.Args[0], "-e", "-o", "command"}

		outBuf, err := psCmd.Output()
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"command": psCmd,
			}).Fatal("Cannot list running processes via ps")
		}

		conflict := false
		for _, batcmd := range r.FindAllString(string(outBuf), -1) {

			log.WithFields(log.Fields{
				"batcmd": batcmd,
			}).Debug("Found a running batsim")

			batargs, err := ParseBatsimCommand(batcmd)
			if err != nil {
				log.WithFields(log.Fields{
					"command": batcmd,
					"err":     err,
				}).Fatal("Cannot retrieve information from a running Batsim process command")
			}

			lineport, err := PortFromBatSock(batargs.Socket)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"extracted socket endpoint": batargs.Socket,
					"batsim command":            batcmd,
				}).Fatal("Cannot retrieve port from a running Batsim process command")
			}

			if lineport == port {
				conflict = true
			}
		}

		if !conflict {
			onexit <- 1
			return
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func waitTcpPortAvailableSs(port uint16, onexit chan int) {
	portStr := strconv.FormatUint(uint64(port), 10)
	r := regexp.MustCompile(":" + portStr)

	for {
		ssCmd := exec.Command("ss")
		ssCmd.Args = []string{ssCmd.Args[0], "-tln"}

		outBuf, err := ssCmd.Output()
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"command": ssCmd,
			}).Fatal("Cannot list open sockets via ss")
		}

		if !(r.Match(outBuf)) {
			onexit <- 1
			return
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func logExecuteTimeoutError(errMsg string, err error,
	name, cmdString, cmdFile, stdoutFile, stderrFile string,
	cmd *exec.Cmd, timeout float64, previewOnError bool) {

	log.WithFields(log.Fields{
		"process name":                 name,
		"err":                          err,
		"command":                      cmdString,
		"command file":                 cmdFile,
		"stdout file":                  stdoutFile,
		"stderr file":                  stderrFile,
		"simulation timeout (seconds)": timeout,
	}).Error(errMsg)

	// If the option is set, preview simulation logs to stdout
	if previewOnError {
		var linesToPreview int64 = 20
		// Preview stdout log unless set to /dev/null (batsim process)
		if stdoutFile != "/dev/null" {
			outPreview, err := PreviewFile(stdoutFile, linesToPreview)
			if err == nil {
				if outPreview != "" {
					fmt.Printf("\nContent of %s's stdout log:\n%s\n",
						name, outPreview)
				}
			} else {
				fmt.Printf("Cannot read %s's stdout log (err=%s)",
					name, err.Error())
			}
		}

		// Preview stderr
		errPreview, err := PreviewFile(stderrFile, linesToPreview)
		if err == nil {
			if errPreview != "" {
				fmt.Printf("\nContent of %s's stderr log:\n%s\n",
					name, errPreview)
			}
		} else {
			fmt.Printf("Cannot read %s's stderr log (err=%s)",
				name, err.Error())
		}
	}
}

// Execute a command, writing status result on a channel
func ExecuteTimeout(name, cmdString, cmdFile, stdoutFile, stderrFile,
	subprocessType string,
	cmd *exec.Cmd, timeout float64, onstart chan CmdFinishedMsg,
	onexit chan CmdFinishedMsg, previewOnError bool) {

	log.WithFields(log.Fields{
		"process name": name,
		"command":      cmdString,
		"command file": cmdFile,
		"timeout":      timeout,
	}).Debug(fmt.Sprintf("Starting %s subprocess", subprocessType))

	if err := cmd.Start(); err != nil {
		// Start failed
		log.WithFields(log.Fields{
			"process name": name,
			"command":      cmdString,
			"command file": cmdFile,
			"stdout file":  stdoutFile,
			"stderr file":  stderrFile,
		}).Error(fmt.Sprintf("Could not start %s subprocess",
			subprocessType))
		onstart <- CmdFinishedMsg{name, FAILURE}
		onexit <- CmdFinishedMsg{name, FAILURE}
		return
	}

	// Start succeeded
	pid := cmd.Process.Pid
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	onstart <- CmdFinishedMsg{name, SUCCESS}

	// Wait until command completion (or context timeout)
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		logExecuteTimeoutError(
			fmt.Sprintf("%s subprocess failed (simulation timeout reached)",
				subprocessType), nil,
			name, cmdString, cmdFile, stdoutFile, stderrFile, cmd, timeout,
			previewOnError)
		KillProcess(pid)
		onexit <- CmdFinishedMsg{name, TIMEOUT}
	case err := <-done:
		if err != nil {
			logExecuteTimeoutError(
				fmt.Sprintf("%s subprocess failed", subprocessType), err,
				name, cmdString, cmdFile, stdoutFile, stderrFile, cmd, timeout,
				previewOnError)
			KillProcess(pid)
			onexit <- CmdFinishedMsg{name, FAILURE}
		} else {
			log.WithFields(log.Fields{
				"process name": name,
				"command":      cmdString,
				"command file": cmdFile,
				"stdout file":  stdoutFile,
				"stderr file":  stderrFile,
			}).Info(fmt.Sprintf("%s subprocess succeeded", subprocessType))
			onexit <- CmdFinishedMsg{name, SUCCESS}
		}
	}
}

// "Batsim" <-> "Scheduler"
func oppName(str string) string {
	if str == "Batsim" {
		return "Scheduler"
	} else {
		return "Batsim"
	}
}

func cleanupSubprocesses(pidsToKill map[string]int) {
	for name, pid := range pidsToKill {
		log.WithFields(log.Fields{
			"name": name,
			"pid":  pid,
		}).Warn("Killing process")
		KillProcess(pid)
	}
}

func setupGuards(pidsToKill *map[string]int) {
	// Guard against ctrl+c
	// Does not seem needed, as SIGTERM is also launched in this case

	// Guard against (polite) kill
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		log.Warn("SIGTERM received. Killing remaining subprocesses.")
		cleanupSubprocesses(*pidsToKill)
		os.Exit(3)
	}()
}

func executeBatsimAlone(exp Experiment, previewOnError bool) int {
	log.WithFields(log.Fields{
		"simulation timeout (seconds)": exp.SimulationTimeout,
		"batsim command":               exp.Batcmd,
		"batsim cmdfile":               exp.OutputDir + "/cmd/batsim.bash",
		"batsim logfile":               exp.OutputDir + "/log/batsim.log",
	}).Info("Starting simulation")

	// Create command
	err := ioutil.WriteFile(exp.OutputDir+"/cmd/batsim.bash", []byte(exp.Batcmd), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/cmd/batsim.bash",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	cmd := exec.Command("bash")
	cmd.Args = []string{cmd.Args[0], "-eux",
		exp.OutputDir + "/cmd/batsim.bash"}

	// Log simulation output
	batlog, err := os.Create(exp.OutputDir + "/log/batsim.log")
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "log/batsim.log",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	defer batlog.Close()
	cmd.Stderr = batlog

	// Guards against SIGINT (ctrl+c) and SIGTERM (polite kill)
	pidsToKill := make(map[string]int)
	setupGuards(&pidsToKill)

	// Execute the processes
	start := make(chan CmdFinishedMsg)
	termination := make(chan CmdFinishedMsg)
	go ExecuteTimeout("Batsim", exp.Batcmd, exp.OutputDir+"/cmd/batsim.bash",
		"/dev/null", exp.OutputDir+"/log/batsim.log", "Simulation", cmd,
		exp.SimulationTimeout, start, termination, previewOnError)

	for {
		select {
		case start1 := <-start:
			if start1.State == SUCCESS {
				pidsToKill["Batsim"] = cmd.Process.Pid
			}
		case finish1 := <-termination:
			delete(pidsToKill, "Batsim")
			return finish1.State
		}
	}
}

func executeBatsimAndSched(exp Experiment, previewOnError bool) int {
	log.WithFields(log.Fields{
		"simulation timeout (seconds)": exp.SimulationTimeout,
		"batsim command":               exp.Batcmd,
		"batsim cmdfile":               exp.OutputDir + "/cmd/batsim.bash",
		"batsim logfile":               exp.OutputDir + "/log/batsim.log",
		"scheduler command":            exp.Schedcmd,
		"scheduler cmdfile":            exp.OutputDir + "/cmd/sched.bash",
		"scheduler logfile (out)":      exp.OutputDir + "/log/sched.out.log",
		"scheduler logfile (err)":      exp.OutputDir + "/log/sched.err.log",
	}).Info("Starting simulation")

	// Create commands
	cmds := make(map[string]*exec.Cmd)
	success := make(map[string]int)

	err := ioutil.WriteFile(exp.OutputDir+"/cmd/batsim.bash", []byte(exp.Batcmd), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/cmd/batsim.bash",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	cmds["Batsim"] = exec.Command("bash")
	cmds["Batsim"].SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // To kill subprocesses later on
	cmds["Batsim"].Args = []string{cmds["Batsim"].Args[0], "-eux",
		exp.OutputDir + "/cmd/batsim.bash"}

	err = ioutil.WriteFile(exp.OutputDir+"/cmd/sched.bash", []byte(exp.Schedcmd), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/cmd/sched.bash",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	cmds["Scheduler"] = exec.Command("bash")
	cmds["Scheduler"].SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // To kill subprocesses later on
	cmds["Scheduler"].Args = []string{cmds["Scheduler"].Args[0], "-eux",
		exp.OutputDir + "/cmd/sched.bash"}

	// Log simulation output
	batlog, err := os.Create(exp.OutputDir + "/log/batsim.log")
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "log/batsim.log",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	defer batlog.Close()
	cmds["Batsim"].Stderr = batlog

	schedout, err := os.Create(exp.OutputDir + "/log/sched.out.log")
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/log/sched.out.log",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	defer schedout.Close()
	cmds["Scheduler"].Stdout = schedout

	schederr, err := os.Create(exp.OutputDir + "/log/sched.err.log")
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/log/sched.err.log",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	defer schederr.Close()
	cmds["Scheduler"].Stderr = schederr

	// Guards against SIGINT (ctrl+c) and SIGTERM (polite kill)
	pidsToKill := make(map[string]int)
	setupGuards(&pidsToKill)

	// Execute the processes
	start := make(chan CmdFinishedMsg)
	termination := make(chan CmdFinishedMsg)
	go ExecuteTimeout("Batsim", exp.Batcmd, exp.OutputDir+"/cmd/batsim.bash",
		"/dev/null", exp.OutputDir+"/log/batsim.log", "Simulation",
		cmds["Batsim"], exp.SimulationTimeout, start, termination,
		previewOnError)
	go ExecuteTimeout("Scheduler", exp.Schedcmd,
		exp.OutputDir+"/cmd/sched.bash",
		exp.OutputDir+"/log/sched.out.log", exp.OutputDir+"/log/sched.err.log",
		"Simulation", cmds["Scheduler"], exp.SimulationTimeout, start,
		termination, previewOnError)

	// Wait for both to start (or to fail starting)
	nbStartedOrFailedStarting := 0
	for nbStartedOrFailedStarting < 2 {
		select {
		case start1 := <-start:
			if start1.State == SUCCESS {
				pidsToKill[start1.Name] = cmds[start1.Name].Process.Pid
			}
			nbStartedOrFailedStarting += 1
		}
	}

	// Wait for first process to finish
	var finish1, finish2 CmdFinishedMsg
	finish1 = <-termination
	delete(pidsToKill, finish1.Name)
	success[finish1.Name] = finish1.State

	log.WithFields(log.Fields{
		"name":  finish1.Name,
		"state": finish1.State,
	}).Debug("First process finished")

	// Depending on the first process success state, we'll wait the second
	// process differently.
	switch finish1.State {
	case SUCCESS:
		log.WithFields(log.Fields{
			"success timeout (seconds)": exp.SuccessTimeout,
			"potential victim name":     oppName(finish1.Name),
		}).Info("The second process might be killed soon...")

		select {
		case <-time.After(time.Duration(exp.SuccessTimeout) * time.Second):
			// Success timeout reached
			log.WithFields(log.Fields{
				"success timeout (seconds)": exp.SuccessTimeout,
			}).Warn("Success timeout reached")

			// Kill the other process
			KillProcess(pidsToKill[oppName(finish1.Name)])
			finish2 = <-termination
		case finish2 = <-termination:
		}
	case FAILURE:
		log.WithFields(log.Fields{
			"failure timeout (seconds)": exp.FailureTimeout,
			"potential victim name":     oppName(finish1.Name),
		}).Info("The second process might be killed soon...")

		select {
		case <-time.After(time.Duration(exp.FailureTimeout) * time.Second):
			// Failure timeout reached
			log.WithFields(log.Fields{
				"failure timeout (seconds)": exp.FailureTimeout,
			}).Warn("Failure timeout reached")

			// Kill the other process
			KillProcess(pidsToKill[oppName(finish1.Name)])
			finish2 = <-termination
		case finish2 = <-termination:
		}
	case TIMEOUT:
		// Wait second process completion
		finish2 = <-termination
	}

	// Second process finished
	delete(pidsToKill, finish2.Name)
	success[finish2.Name] = finish2.State

	return max(success["Batsim"], success["Scheduler"])
}

// Execute one Batsim simulation
func ExecuteOne(exp Experiment, previewOnError bool) int {
	// Prepare execution
	PrepareDirs(exp)

	// Sets unset command as empty string
	if exp.Schedcmd == "schedcmd-unset" {
		exp.Schedcmd = ""
	}

	// Parse batsim command
	batargs, err := ParseBatsimCommand(exp.Batcmd)
	if err != nil {
		log.WithFields(log.Fields{
			"command": exp.Batcmd,
			"err":     err,
		}).Fatal("Cannot retrieve information from Batsim command")
	}

	if !strings.HasPrefix(batargs.ExportPrefix, exp.OutputDir) {
		log.WithFields(log.Fields{
			"batsim prefix":    batargs.ExportPrefix,
			"output directory": exp.OutputDir,
			"batsim command":   exp.Batcmd,
		}).Warning("Batsim export prefix mismatches output directory")
	}

	if exp.Schedcmd == "" {
		// Only execute Batsim
		if batargs.BatexecMode == false {
			log.WithFields(log.Fields{
				"batsim command":    exp.Batcmd,
				"scheduler command": exp.Schedcmd,
			}).Fatal("Sched command unset but Batsim is not in batexec mode")
		}

		return executeBatsimAlone(exp, previewOnError)
	} else {
		// Execute Batsim and the scheduler
		if batargs.BatexecMode == true {
			log.WithFields(log.Fields{
				"batsim command":    exp.Batcmd,
				"scheduler command": exp.Schedcmd,
			}).Fatal("Sched command set but Batsim is in batexec mode")
		}

		// Wait for context to be ready (open sockets, batsim processes...)
		waitReadyForSimulation(exp, batargs)

		return executeBatsimAndSched(exp, previewOnError)
	}

	return -1
}

func KillProcess(pid int) {
	syscall.Kill(-pid, syscall.SIGTERM)
}
