package main

import (
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
	EXPECT_ABSENCE
)

const (
	UNEXPECTED_ROBIN_SUCCESS_STATE int = 1 << iota
	UNEXPECTED_BATSIM_SUCCESS_STATE
	UNEXPECTED_SCHED_SUCCESS_STATE
	UNEXPECTED_SCHED_PRESENCE_STATE
	UNEXPECTED_CONTEXT_CLEANLINESS
)

func main() {
	usage := `Tests one robin execution.

Usage: 
  robintest <description-file>
  			--test-timeout=<seconds>
  			[(--expect-robin-success | --expect-robin-failure)]
  			[(--expect-batsim-success | --expect-batsim-failure)]
  			[(--expect-sched-success | --expect-sched-failure | --expect-no-sched)]
  			[(--expect-ctx-clean | --expect-ctx-busy)]
  			[(--expect-ctx-clean-at-begin | --expect-ctx-busy-at-begin)]
  			[(--expect-ctx-clean-at-end | --expect-ctx-busy-at-end)]

  robintest -h | --help
  robintest --version`

	arguments, _ := docopt.Parse(usage, nil, true, batexpe.Version(), false)

	robinExpectation := EXPECT_NOTHING
	if arguments["--expect-robin-success"] == true {
		robinExpectation = EXPECT_TRUE
	} else if arguments["--expect-robin-failure"] == true {
		robinExpectation = EXPECT_FALSE
	}

	ctxExpectation := EXPECT_NOTHING
	if arguments["--expect-ctx-clean"] == true {
		ctxExpectation = EXPECT_TRUE
	} else if arguments["--expect-ctx-busy"] == true {
		ctxExpectation = EXPECT_FALSE
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
	} else if arguments["--expect-no-sched"] == true {
		schedExpectation = EXPECT_ABSENCE
	}

	testTimeout, err := strconv.ParseFloat(arguments["--test-timeout"].(string), 64)
	if err != nil {
		panic("Invalid test-timeout value: Cannot convert to float")
	}

	testResult := RobinTest(arguments["<description-file>"].(string),
		robinExpectation, batsimExpectation, schedExpectation,
		ctxExpectation, ctxExpectationAtEnd, testTimeout)

	os.Exit(testResult)
}

func RobinTest(descriptionFile string,
	robinExpectation, batsimExpectation, schedExpectation,
	ctxExpectation, ctxExpectationAtEnd int,
	testTimeout float64) int {

	returnValue := 0

	rresult := batexpe.RunRobin(descriptionFile, testTimeout)

	jsonLines, err := batexpe.ParseRobinOutput(rresult.Output)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not parse robin output")
	}

	// Robin result
	if robinExpectation != EXPECT_NOTHING {
		expectedRobinSuccess := robinExpectation == EXPECT_TRUE
		robinSuccess := rresult.Finished && rresult.ReturnCode == 0

		if robinSuccess != expectedRobinSuccess {
			log.WithFields(log.Fields{
				"expected": expectedRobinSuccess,
				"got":      robinSuccess,
			}).Error("Unexpected robin success state")

			returnValue += UNEXPECTED_ROBIN_SUCCESS_STATE
		}
	}

	// Batsim successfulness
	if batsimExpectation != EXPECT_NOTHING {
		expectedBatsimSuccess := batsimExpectation == EXPECT_TRUE
		batsimSuccess := batexpe.WasBatsimSuccessful(jsonLines)

		if batsimSuccess != expectedBatsimSuccess {
			log.WithFields(log.Fields{
				"expected": expectedBatsimSuccess,
				"got":      batsimSuccess,
			}).Error("Unexpected batsim success state")

			returnValue += UNEXPECTED_BATSIM_SUCCESS_STATE
		}
	}

	// Sched successfulness and presence
	if schedExpectation != EXPECT_NOTHING {
		expectedSchedSuccess := schedExpectation == EXPECT_TRUE
		expectedSchedPresence := schedExpectation != EXPECT_ABSENCE
		schedSuccess, schedPresence := batexpe.WasSchedSuccessful(jsonLines)

		if schedSuccess != expectedSchedSuccess {
			log.WithFields(log.Fields{
				"expected": expectedSchedSuccess,
				"got":      schedSuccess,
			}).Error("Unexpected sched success state")

			returnValue += UNEXPECTED_SCHED_SUCCESS_STATE
		}

		if schedPresence != expectedSchedPresence {
			log.WithFields(log.Fields{
				"expected": expectedSchedPresence,
				"got":      schedPresence,
			}).Error("Unexpected sched presence state")

			returnValue += UNEXPECTED_SCHED_PRESENCE_STATE
		}
	}

	// Context cleanliness during simulation
	if ctxExpectation != EXPECT_NOTHING {
		expectedContextClean := ctxExpectation == EXPECT_TRUE
		contextClean := batexpe.WasContextClean(jsonLines)

		if contextClean != expectedContextClean {
			log.WithFields(log.Fields{
				"expected": expectedContextClean,
				"got":      contextClean,
			}).Error("Unexpected context cleanliness")

			returnValue += UNEXPECTED_CONTEXT_CLEANLINESS
		}
	}

	return returnValue
}
