package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"framagit.org/batsim/batexpe"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var (
	version string
)

const (
	EXPECT_NOTHING int = iota
	EXPECT_TRUE
	EXPECT_FALSE
	EXPECT_ABSENCE
	EXPECT_KILLED
)

func setupLogging(arguments map[string]interface{}) {
	log.SetOutput(os.Stdout)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
	customFormatter.FullTimestamp = true
	customFormatter.QuoteEmptyFields = true
	log.SetFormatter(customFormatter)

	if arguments["--debug"] == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	os.Exit(mainReturnWithCode())
}

func mainReturnWithCode() int {
	usage := `Tests one robin execution.

Usage: 
  robintest <description-file>
  			--test-timeout=<seconds>
  			[(--expect-robin-success | --expect-robin-failure |
  			  --expect-robin-killed)]
  			[(--expect-batsim-success | --expect-batsim-failure |
  			  --expect-batsim-killed)]
  			[(--expect-sched-success | --expect-sched-failure |
  			  --expect-sched-killed | --expect-no-sched)]
  			[(--expect-ctx-clean | --expect-ctx-busy)]
  			[(--expect-ctx-clean-at-begin | --expect-ctx-busy-at-begin)]
  			[(--expect-ctx-clean-at-end | --expect-ctx-busy-at-end)]
  			[--result-check-script=<file>]
  			[--cover=<file>]
  			[--debug]
  robintest -h | --help
  robintest --version`

	robintestVersion := version
	if robintestVersion == "" {
		robintestVersion = batexpe.Version()
	}

	ret := -1

	parser := &docopt.Parser{
		HelpHandler: func(err error, usage string) {
			fmt.Println(usage)
			if err != nil {
				ret = 1
			} else {
				ret = 0
			}
		},
		OptionsFirst: false,
	}

	arguments, _ := parser.ParseArgs(usage, os.Args[1:], robintestVersion)
	if ret != -1 {
		return ret
	}

	setupLogging(arguments)

	// Has robin been successful? (returned 0 before test timeout)
	robinExpectation := EXPECT_NOTHING
	if arguments["--expect-robin-success"] == true {
		robinExpectation = EXPECT_TRUE
	} else if arguments["--expect-robin-failure"] == true {
		robinExpectation = EXPECT_FALSE
	} else if arguments["--expect-robin-killed"] == true {
		robinExpectation = EXPECT_KILLED
	}

	// Did the execution context become clean during robin's execution?
	ctxExpectation := EXPECT_NOTHING
	if arguments["--expect-ctx-clean"] == true {
		ctxExpectation = EXPECT_TRUE
	} else if arguments["--expect-ctx-busy"] == true {
		ctxExpectation = EXPECT_FALSE
	}

	// Was the execution context clean before running robin?
	ctxExpectationAtBegin := EXPECT_NOTHING
	if arguments["--expect-ctx-clean-at-begin"] == true {
		ctxExpectationAtBegin = EXPECT_TRUE
	} else if arguments["--expect-ctx-busy-at-begin"] == true {
		ctxExpectationAtBegin = EXPECT_FALSE
	}

	// Was the execution context clean after running robin?
	ctxExpectationAtEnd := EXPECT_NOTHING
	if arguments["--expect-ctx-clean-at-end"] == true {
		ctxExpectationAtEnd = EXPECT_TRUE
	} else if arguments["--expect-ctx-busy-at-end"] == true {
		ctxExpectationAtEnd = EXPECT_FALSE
	}

	// Was Batsim successful in robin's execution?
	batsimExpectation := EXPECT_NOTHING
	if arguments["--expect-batsim-success"] == true {
		batsimExpectation = EXPECT_TRUE
	} else if arguments["--expect-batsim-failure"] == true {
		batsimExpectation = EXPECT_FALSE
	} else if arguments["--expect-batsim-killed"] == true {
		batsimExpectation = EXPECT_KILLED
	}

	// Was the scheduler successful (and present) in robin's execution?
	schedExpectation := EXPECT_NOTHING
	if arguments["--expect-sched-success"] == true {
		schedExpectation = EXPECT_TRUE
	} else if arguments["--expect-sched-failure"] == true {
		schedExpectation = EXPECT_FALSE
	} else if arguments["--expect-sched-killed"] == true {
		schedExpectation = EXPECT_KILLED
	} else if arguments["--expect-no-sched"] == true {
		schedExpectation = EXPECT_ABSENCE
	}

	testTimeout, err := strconv.ParseFloat(arguments["--test-timeout"].(string), 64)
	if err != nil {
		log.WithFields(log.Fields{
			"err":            err,
			"--test-timeout": arguments["--test-timeout"].(string),
		}).Error("Invalid test timeout")
		return 1
	}

	coverFile := ""
	if arguments["--cover"] != nil {
		coverFile = arguments["--cover"].(string)
	}

	resultCheckScript := ""
	if arguments["--result-check-script"] != nil {
		resultCheckScript = arguments["--result-check-script"].(string)
	}

	testResult := RobinTest(arguments["<description-file>"].(string),
		coverFile, resultCheckScript, testTimeout,
		robinExpectation, batsimExpectation, schedExpectation,
		ctxExpectation, ctxExpectationAtBegin, ctxExpectationAtEnd)

	return testResult
}

func RobinTest(descriptionFile, coverFile, resultCheckScript string,
	testTimeout float64,
	robinExpectation, batsimExpectation, schedExpectation,
	ctxExpectation, ctxExpectationAtBegin, ctxExpectationAtEnd int) int {

	// Computing whether the context is clean or not is done by checking whether
	// any batsim or batsched is running. This is intentionally done with a
	// different function (and technique) that the one done within robin.

	batRunningAtBegin, err1 := batexpe.IsBatsimOrBatschedRunning()
	ctxCleanAtBegin := batRunningAtBegin == false

	rresult := batexpe.RunRobin(descriptionFile, coverFile, testTimeout)

	batRunningAtEnd, err2 := batexpe.IsBatsimOrBatschedRunning()
	ctxCleanAtEnd := batRunningAtEnd == false

	robintestReturnValue := 0

	jsonLines, parseRobinOutputErr := batexpe.ParseRobinOutput(rresult.Output)

	if (err1 != nil) || (err2 != nil) {
		robintestReturnValue = 1
	}

	// Robin result
	if robinExpectation != EXPECT_NOTHING {
		expectedRobinSuccess := robinExpectation == EXPECT_TRUE
		expectedRobinKilled := robinExpectation == EXPECT_KILLED
		robinSuccess := rresult.Succeeded
		robinKilled := rresult.Finished == false

		if robinSuccess != expectedRobinSuccess {
			log.WithFields(log.Fields{
				"expected": expectedRobinSuccess,
				"got":      robinSuccess,
			}).Error("Unexpected robin success state")

			robintestReturnValue = 1
		}

		if robinKilled != expectedRobinKilled {
			log.WithFields(log.Fields{
				"expected": expectedRobinKilled,
				"got":      robinKilled,
			}).Error("Unexpected robin kill state")

			robintestReturnValue = 1
		}
	}

	// Batsim successfulness
	if batsimExpectation != EXPECT_NOTHING {
		expectedBatsimSuccess := batsimExpectation == EXPECT_TRUE
		expectedBatsimKilled := batsimExpectation == EXPECT_KILLED
		batsimSuccess, batsimKilled := batexpe.WasBatsimSuccessful(jsonLines)

		if batsimSuccess != expectedBatsimSuccess {
			log.WithFields(log.Fields{
				"expected": expectedBatsimSuccess,
				"got":      batsimSuccess,
			}).Error("Unexpected batsim success state")

			robintestReturnValue = 1
		}

		if batsimKilled != expectedBatsimKilled {
			log.WithFields(log.Fields{
				"expected": expectedBatsimKilled,
				"got":      batsimKilled,
			}).Error("Unexpected batsim kill state")

			robintestReturnValue = 1
		}
	}

	// Sched successfulness and presence
	if schedExpectation != EXPECT_NOTHING {
		expectedSchedSuccess := schedExpectation == EXPECT_TRUE
		expectedSchedPresence := schedExpectation != EXPECT_ABSENCE
		expectedSchedKilled := schedExpectation == EXPECT_KILLED
		schedSuccess, schedPresence, schedKilled := batexpe.WasSchedSuccessful(jsonLines)

		if schedSuccess != expectedSchedSuccess {
			log.WithFields(log.Fields{
				"expected": expectedSchedSuccess,
				"got":      schedSuccess,
			}).Error("Unexpected sched success state")

			robintestReturnValue = 1
		}

		if schedPresence != expectedSchedPresence {
			log.WithFields(log.Fields{
				"expected": expectedSchedPresence,
				"got":      schedPresence,
			}).Error("Unexpected sched presence state")

			robintestReturnValue = 1
		}

		if schedKilled != expectedSchedKilled {
			log.WithFields(log.Fields{
				"expected": expectedSchedKilled,
				"got":      schedKilled,
			}).Error("Unexpected sched kill state")

			robintestReturnValue = 1
		}
	}

	// Context cleanliness during robin's execution
	if ctxExpectation != EXPECT_NOTHING {
		expectedContextClean := ctxExpectation == EXPECT_TRUE
		contextClean := batexpe.WasContextClean(jsonLines)

		if contextClean != expectedContextClean {
			log.WithFields(log.Fields{
				"expected": expectedContextClean,
				"got":      contextClean,
			}).Error("Unexpected context cleanliness during robin's execution")

			robintestReturnValue = 1
		}
	}

	// Context cleanliness before robin's execution
	if ctxExpectationAtBegin != EXPECT_NOTHING {
		expectedCtxCleanAtBegin := ctxExpectationAtBegin == EXPECT_TRUE

		if ctxCleanAtBegin != expectedCtxCleanAtBegin {
			log.WithFields(log.Fields{
				"expected": expectedCtxCleanAtBegin,
				"got":      ctxCleanAtBegin,
			}).Error("Unexpected context cleanliness before robin's execution")

			robintestReturnValue = 1
		}
	}

	// Context cleanliness after robin's execution
	if ctxExpectationAtEnd != EXPECT_NOTHING {
		expectedCtxCleanAtEnd := ctxExpectationAtEnd == EXPECT_TRUE

		if ctxCleanAtEnd != expectedCtxCleanAtEnd {
			log.WithFields(log.Fields{
				"expected": expectedCtxCleanAtEnd,
				"got":      ctxCleanAtEnd,
			}).Error("Unexpected context cleanliness after robin's execution")

			robintestReturnValue = 1
		}
	}

	// Run check script if everything went as expected so far
	if robintestReturnValue == 0 {
		if resultCheckScript != "" {
			// First, we need to retrieve Batsim output prefix.
			// To do so, we can parse the batsim command defined in
			// the robin description file.
			byt, err := ioutil.ReadFile(descriptionFile)
			if err != nil {
				log.WithFields(log.Fields{
					"err":      err,
					"filename": descriptionFile,
				}).Error("Cannot open description file")
				robintestReturnValue = 1
			}

			exp, err := batexpe.FromYaml(string(byt))
			if err != nil {
				robintestReturnValue = 1
			}

			batargs, err := batexpe.ParseBatsimCommand(exp.Batcmd)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot parse Batsim command")
				robintestReturnValue = 1
			}

			checkScriptSuccessful, err := RunCheckScript(resultCheckScript,
				exp.OutputDir, batargs.ExportPrefix, testTimeout)
			if err != nil {
				robintestReturnValue = 1
			}

			if !checkScriptSuccessful {
				robintestReturnValue = 1
			}
		}
	}

	if parseRobinOutputErr != nil {
		log.WithFields(log.Fields{
			"err": parseRobinOutputErr,
		}).Error("Could not parse robin output")
		robintestReturnValue = 1
	}

	return robintestReturnValue
}

func RunCheckScript(resultCheckScript, robinOutputDir, batsimExportPrefix string,
	checkTimeout float64) (bool, error) {
	cmd := exec.Command(resultCheckScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // To kill subprocesses later on
	cmd.Args = []string{cmd.Args[0], batsimExportPrefix}

	checkout, outErr := os.Create(robinOutputDir + "/log/check.out.log")
	checkerr, errErr := os.Create(robinOutputDir + "/log/check.err.log")

	if outErr == nil {
		defer checkout.Close()
		cmd.Stdout = checkout
	}

	if errErr == nil {
		defer checkerr.Close()
		cmd.Stderr = checkerr
	}

	if (outErr != nil) || (errErr != nil) {
		log.WithFields(log.Fields{
			"out_filename": robinOutputDir + "/log/check.out.log",
			"out_err":      outErr,
			"err_filename": robinOutputDir + "/log/check.err.log",
			"err_err":      errErr,
		}).Error("Cannot create file")
		return false, fmt.Errorf("Cannot create file")
	}

	start := make(chan batexpe.CmdFinishedMsg)
	termination := make(chan batexpe.CmdFinishedMsg)

	go batexpe.ExecuteTimeout("Check", strings.Join(cmd.Args, " "),
		"/dev/null",
		robinOutputDir+"/log/check.out.log",
		robinOutputDir+"/log/check.err.log",
		"Check", cmd, checkTimeout, start, termination, true)

	<-start
	finish1 := <-termination
	return finish1.State == batexpe.SUCCESS, nil
}
