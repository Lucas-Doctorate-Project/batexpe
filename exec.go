package expe

import (
	"context"
	"github.com/anmitsu/go-shlex"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

const (
	SUCCESS int = iota
	FAILURE
	TIMEOUT
)

type cmdFinishedMsg struct {
	name  string
	state int
}

func PrepareOutput(exp Experiment) {
	// Create output directory if needed
	CreateDirIfNeeded(exp.OutputDir)
	CreateDirIfNeeded(exp.OutputDir + "/log")
}

func executeInnerCtx(name string, cmdString string,
	cmd *exec.Cmd, ctx context.Context, onexit chan cmdFinishedMsg) {

	log.WithFields(log.Fields{
		"command": cmdString,
		"context": ctx,
	}).Debug("Running " + name)

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
	} else {
		log.WithFields(log.Fields{
			"command": cmdString,
		}).Info(name + " execution succeeded")

		onexit <- cmdFinishedMsg{name, SUCCESS}
	}
}

func oppName(str string) string {
	if str == "Batsim" {
		return "Scheduler"
	} else {
		return "Batsim"
	}
}

func ExecuteOne(exp Experiment) int {
	// Prepare execution
	PrepareOutput(exp)

	// Sets unset command as empty string
	if exp.Schedcmd == "schedcmd-unset" {
		exp.Schedcmd = ""
	}

	// Parse batsim command
	batargs := ParseBatsimCommand(exp.Batcmd)

	// TODO: check whether export prefix matches output dir

	// Only execute Batsim
	if exp.Schedcmd == "" {
		if batargs.BatexecMode == false {
			log.WithFields(log.Fields{
				"batsim command":    exp.Batcmd,
				"scheduler command": exp.Schedcmd,
			}).Fatal("Sched command unset but Batsim is not in batexec mode")
		}

		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(exp.SimulationTimeout)*time.Second)
		defer cancel()

		log.WithFields(log.Fields{
			"batsim command": exp.Batcmd,
			"timeout":        exp.SimulationTimeout,
		}).Info("Running Batsim")

		splitBatsimCmd, err := shlex.Split(exp.Batcmd, true)
		if err != nil {
			log.WithFields(log.Fields{
				"batsim command": exp.Batcmd,
				"err":            err,
			}).Fatal("Shlex split failed")
		}
		cmd := exec.CommandContext(ctx, splitBatsimCmd[0])
		cmd.Args = splitBatsimCmd

		termination := make(chan cmdFinishedMsg)
		go executeInnerCtx("Batsim", exp.Batcmd, cmd, ctx, termination)
		finish1 := <-termination
		return finish1.state
	} else {
		// Execute Batsim and the scheduler
		if batargs.BatexecMode == true {
			log.WithFields(log.Fields{
				"batsim command":    exp.Batcmd,
				"scheduler command": exp.Schedcmd,
			}).Fatal("Sched command set but Batsim is in batexec mode")
		}

		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(exp.SimulationTimeout)*time.Second)
		defer cancel()

		log.WithFields(log.Fields{
			"batsim command":    exp.Batcmd,
			"scheduler command": exp.Schedcmd,
			"timeout (seconds)": exp.SimulationTimeout,
		}).Info("Running Batsim and Scheduler")

		cmds := make(map[string]*exec.Cmd)
		success := make(map[string]int)

		splitBatsimCmd, err := shlex.Split(exp.Batcmd, true)
		if err != nil {
			log.WithFields(log.Fields{
				"batsim command": exp.Batcmd,
				"err":            err,
			}).Fatal("Shlex split failed")
		}
		cmds["Batsim"] = exec.CommandContext(ctx, splitBatsimCmd[0])
		cmds["Batsim"].Args = splitBatsimCmd

		splitSchedCmd, err := shlex.Split(exp.Schedcmd, true)
		if err != nil {
			log.WithFields(log.Fields{
				"scheduler command": exp.Schedcmd,
				"err":               err,
			}).Fatal("Shlex split failed")
		}
		cmds["Scheduler"] = exec.CommandContext(ctx, splitSchedCmd[0])
		cmds["Scheduler"].Args = splitSchedCmd

		termination := make(chan cmdFinishedMsg)
		go executeInnerCtx("Batsim", exp.Batcmd, cmds["Batsim"], ctx, termination)
		go executeInnerCtx("Scheduler", exp.Schedcmd, cmds["Scheduler"], ctx, termination)

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

	return -1
}
