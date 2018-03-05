package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"gitlab.inria.fr/batsim/batexpe"
	"os"
	"strconv"
)

const (
	EXPECT_NOTHING int = iota
	EXPECT_TRUE
	EXPECT_FALSE
)

func main() {
	usage := `Tests one robin execution.

Usage: 
  robintest <description-file>
  			[(--expect-simu-success | --expect-simu-failure)]
  			[(--expect-batsim-success | --expect-batsim-failure)]
  			[(--expect-sched-success | --expect-sched-failure)]
  			[(--expect-ctx-clean | --expect-ctx-busy)]
  			[(--expect-ctx-clean-at-end | --expect-ctx-busy-at-end)]
  			--test-timeout=<seconds>

  robintest -h | --help
  robintest --version`

	arguments, _ := docopt.Parse(usage, nil, true, batexpe.Version(), false)
	fmt.Print(arguments)

	simuExpectation := EXPECT_NOTHING
	if arguments["--expect-simu-success"] == true {
		simuExpectation = EXPECT_TRUE
	} else if arguments["--expect-simu-failure"] == true {
		simuExpectation = EXPECT_FALSE
	}

	ctxExpectation := EXPECT_NOTHING
	if arguments["--expect-ctx-clean"] == true {
		ctxExpectation = EXPECT_TRUE
	} else if arguments["--expect-ctx-busy"] == true {
		simuExpectation = EXPECT_FALSE
	}

	ctxExpectationAtEnd := EXPECT_NOTHING
	if arguments["--expect-ctx-clean-at-end"] == true {
		ctxExpectationAtEnd = EXPECT_TRUE
	} else if arguments["--expect-ctx-busy-at-end"] == true {
		ctxExpectationAtEnd = EXPECT_FALSE
	}

	batsimExpectation := EXPECT_NOTHING
	if arguments["--expect-batsim-success"] == true {
		batsimExpectation = EXPECT_TRUE
	} else if arguments["--expect-batsim-failure"] == true {
		batsimExpectation = EXPECT_FALSE
	}

	schedExpectation := EXPECT_NOTHING
	if arguments["--expect-sched-success"] == true {
		schedExpectation = EXPECT_TRUE
	} else if arguments["--expect-sched-failure"] == true {
		schedExpectation = EXPECT_FALSE
	}

	testTimeout, err := strconv.ParseFloat(arguments["--test-timeout"].(string), 64)
	if err != nil {
		panic("Invalid test-timeout value: Cannot convert to float")
	}

	testResult := RobinTest(arguments["<description-file>"].(string),
		simuExpectation, batsimExpectation, schedExpectation,
		ctxExpectation, ctxExpectationAtEnd, testTimeout)

	os.Exit(testResult)
}

func RobinTest(descriptionFile string,
	simuExpectation, batsimExpectation, schedExpectation,
	ctxExpectation, ctxExpectationAtEnd int,
	testTimeout float64) int {

	rresult := batexpe.RunRobin(descriptionFile, testTimeout)

	log.WithFields(log.Fields{
		"output": rresult.Output,
	}).Info("Retrieved robin output")

	json_lines, err := batexpe.ParseRobinOutput(rresult.Output)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not parse robin output")
	}

	log.WithFields(log.Fields{
		"lines": json_lines,
	}).Info("Parsed robin output")

	return 0
}
