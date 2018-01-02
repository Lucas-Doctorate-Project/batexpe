package expe

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	SUCCESS int = iota
	TIMEOUT
	FAILURE
)

type cmdFinishedMsg struct {
	name  string
	state int
}

func PrepareDirs(exp Experiment) {
	// Create output directory if needed
	CreateDirIfNeeded(exp.OutputDir)
	CreateDirIfNeeded(exp.OutputDir + "/log")
	CreateDirIfNeeded(exp.OutputDir + "/cmd")
}

// Execute a command, writing status result on a channel
func executeInnerCtx(name string, cmdString string,
	cmd *exec.Cmd, ctx context.Context, onexit chan cmdFinishedMsg) {

	log.WithFields(log.Fields{
		"command": cmdString,
		"context": ctx,
	}).Debug("Starting " + name)

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"command": cmdString,
			}).Error(name + " execution failed (simulation timeout reached)")

			onexit <- cmdFinishedMsg{name, TIMEOUT}
		} else {
			log.WithFields(log.Fields{
				"err":     err,
				"command": cmdString,
			}).Error(name + " execution failed")

			onexit <- cmdFinishedMsg{name, FAILURE}
		}

		// Kill command's subprocesses (cf. https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	} else {
		log.WithFields(log.Fields{
			"command": cmdString,
		}).Info(name + " execution succeeded")

		onexit <- cmdFinishedMsg{name, SUCCESS}
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

func executeBatsimAlone(exp Experiment, ctx context.Context) int {
	log.WithFields(log.Fields{
		"simulation timeout (seconds)": exp.SimulationTimeout,
		"batsim command":               exp.Batcmd,
		"batsim cmdfile":               exp.OutputDir + "/cmd/batsim.bash",
		"batsim logfile":               exp.OutputDir + "/log/batsim.log",
	}).Info("Running Batsim")

	// Create command
	err := ioutil.WriteFile(exp.OutputDir+"/cmd/batsim.bash", []byte(exp.Batcmd), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/cmd/batsim.bash",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	cmd := exec.CommandContext(ctx, "bash")
	cmd.Args = []string{"bash", "-eux", exp.OutputDir + "/cmd/batsim.bash"}

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

	// Execute the processes
	termination := make(chan cmdFinishedMsg)
	go executeInnerCtx("Batsim", exp.Batcmd, cmd, ctx, termination)

	// Guard against ctrl+c
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint

		log.WithFields(log.Fields{
			"Batsim cmd": cmd,
		}).Warn("SIGTERM received. Killing subprocesses.")

		syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
		os.Exit(3)
	}()

	// Wait for completion
	finish1 := <-termination
	return finish1.state
}

func executeBatsimAndSched(exp Experiment, ctx context.Context) int {
	log.WithFields(log.Fields{
		"simulation timeout (seconds)": exp.SimulationTimeout,
		"batsim command":               exp.Batcmd,
		"batsim cmdfile":               exp.OutputDir + "/cmd/batsim.bash",
		"batsim logfile":               exp.OutputDir + "/log/batsim.log",
		"scheduler command":            exp.Schedcmd,
		"scheduler cmdfile":            exp.OutputDir + "/cmd/sched.bash",
		"scheduler logfile (out)":      exp.OutputDir + "/log/sched.out.log",
		"scheduler logfile (err)":      exp.OutputDir + "/log/sched.err.log",
	}).Info("Running Batsim and Scheduler")

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
	cmds["Batsim"] = exec.CommandContext(ctx, "bash")
	cmds["Batsim"].SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // To kill subprocesses later on
	cmds["Batsim"].Args = []string{"bash", "-eux", exp.OutputDir + "/cmd/batsim.bash"}

	err = ioutil.WriteFile(exp.OutputDir+"/cmd/sched.bash", []byte(exp.Schedcmd), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": exp.OutputDir + "/cmd/sched.bash",
			"err":      err,
		}).Fatal("Cannot create file")
	}
	cmds["Scheduler"] = exec.CommandContext(ctx, "bash")
	cmds["Scheduler"].SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // To kill subprocesses later on
	cmds["Scheduler"].Args = []string{"bash", "-eux", exp.OutputDir + "/cmd/sched.bash"}

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

	// Execute the processes
	termination := make(chan cmdFinishedMsg)
	go executeInnerCtx("Batsim", exp.Batcmd, cmds["Batsim"], ctx, termination)
	go executeInnerCtx("Scheduler", exp.Schedcmd, cmds["Scheduler"], ctx, termination)

	// Guard against ctrl+c
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint

		log.WithFields(log.Fields{
			"Batsim cmd":    cmds["Batsim"],
			"Scheduler cmd": cmds["Scheduler"],
		}).Warn("SIGTERM received. Killing subprocesses.")

		syscall.Kill(-cmds["Batsim"].Process.Pid, syscall.SIGKILL)
		syscall.Kill(-cmds["Scheduler"].Process.Pid, syscall.SIGKILL)
		os.Exit(3)
	}()

	// Wait for first process to finish
	finish1 := <-termination
	success[finish1.name] = finish1.state

	log.WithFields(log.Fields{
		"name":  finish1.name,
		"state": finish1.state,
	}).Debug("First process finished")

	switch finish1.state {
	case SUCCESS:
		log.WithFields(log.Fields{
			"success timeout (seconds)": exp.SuccessTimeout,
			"potential victim":          oppName(finish1.name),
		}).Info(oppName(finish1.name) + " may be killed soon...")

		select {
		case <-time.After(time.Duration(exp.SuccessTimeout) * time.Second):
			// Success timeout reached
			log.WithFields(log.Fields{
				"success timeout (seconds)": exp.SuccessTimeout,
			}).Warn("Success timeout reached")

			// Kill the other process
			if err := cmds[oppName(finish1.name)].Process.Kill(); err != nil {
				log.WithFields(log.Fields{
					"name": oppName(finish1.name),
					"err":  err,
				}).Error("Failed to kill")
			}

			// Wait second process completion
			finish2 := <-termination
			success[finish2.name] = finish2.state

		case finish2 := <-termination:
			// Second process finished
			success[finish2.name] = finish2.state
		}
	case FAILURE:
		log.WithFields(log.Fields{
			"failure timeout (seconds)": exp.FailureTimeout,
			"potential victim":          oppName(finish1.name),
		}).Info(oppName(finish1.name) + " may be killed soon...")

		select {
		case <-time.After(time.Duration(exp.FailureTimeout) * time.Second):
			// Failure timeout reached
			log.WithFields(log.Fields{
				"failure timeout (seconds)": exp.FailureTimeout,
			}).Warn("Failure timeout reached")

			// Kill the other process
			if err := cmds[oppName(finish1.name)].Process.Kill(); err != nil {
				log.WithFields(log.Fields{
					"name": oppName(finish1.name),
					"err":  err,
				}).Error("Failed to kill")
			}

			// Wait second process completion
			finish2 := <-termination
			success[finish2.name] = finish2.state

		case finish2 := <-termination:
			// Second process finished
			success[finish2.name] = finish2.state
		}
	case TIMEOUT:
		// Wait second process completion
		finish2 := <-termination
		success[finish2.name] = finish2.state
	}

	return max(success["Batsim"], success["Scheduler"])
}

// Execute one Batsim simulation
func ExecuteOne(exp Experiment) int {
	// Prepare execution
	PrepareDirs(exp)

	// Sets unset command as empty string
	if exp.Schedcmd == "schedcmd-unset" {
		exp.Schedcmd = ""
	}

	// Parse batsim command
	batargs := ParseBatsimCommand(exp.Batcmd)

	if !strings.HasPrefix(batargs.ExportPrefix, exp.OutputDir) {
		log.WithFields(log.Fields{
			"batsim prefix":    batargs.ExportPrefix,
			"output directory": exp.OutputDir,
			"batsim command":   exp.Batcmd,
		}).Warning("Batsim export prefix mismatches output directory")
	}

	// Create context (handles timeout)
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(exp.SimulationTimeout)*time.Second)
	defer cancel()

	if exp.Schedcmd == "" {
		// Only execute Batsim
		if batargs.BatexecMode == false {
			log.WithFields(log.Fields{
				"batsim command":    exp.Batcmd,
				"scheduler command": exp.Schedcmd,
			}).Fatal("Sched command unset but Batsim is not in batexec mode")
		}
		return executeBatsimAlone(exp, ctx)
	} else {
		// Execute Batsim and the scheduler
		if batargs.BatexecMode == true {
			log.WithFields(log.Fields{
				"batsim command":    exp.Batcmd,
				"scheduler command": exp.Schedcmd,
			}).Fatal("Sched command set but Batsim is in batexec mode")
		}
		return executeBatsimAndSched(exp, ctx)
	}

	return -1
}
