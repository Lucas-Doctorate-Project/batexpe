package main

import (
	docopt "github.com/docopt/docopt-go"
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
	EXPECT_KILLED
)

const (
	UNEXPECTED_ROBIN_SUCCESS_STATE int = 1 << iota
	UNEXPECTED_ROBIN_KILL_STATE
	UNEXPECTED_BATSIM_SUCCESS_STATE
	UNEXPECTED_BATSIM_KILL_STATE
	UNEXPECTED_SCHED_SUCCESS_STATE
	UNEXPECTED_SCHED_PRESENCE_STATE
	UNEXPECTED_SCHED_KILL_STATE
	UNEXPECTED_CONTEXT_CLEANLINESS
	UNEXPECTED_CONTEXT_CLEANLINESS_AT_BEGIN
	UNEXPECTED_CONTEXT_CLEANLINESS_AT_END
)

func main() {
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
  			[--cover=<file>]

  robintest -h | --help
  robintest --version`

	arguments, _ := docopt.Parse(usage, nil, true, batexpe.Version(), false)

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
		panic("Invalid test-timeout value: Cannot convert to float")
	}

	coverFile := ""
	if arguments["--cover"] != nil {
		coverFile = arguments["--cover"].(string)
	}

	testResult := RobinTest(arguments["<description-file>"].(string),
		coverFile, testTimeout,
		robinExpectation, batsimExpectation, schedExpectation,
		ctxExpectation, ctxExpectationAtBegin, ctxExpectationAtEnd)

	os.Exit(testResult)
}

func RobinTest(descriptionFile, coverFile string, testTimeout float64,
	robinExpectation, batsimExpectation, schedExpectation,
	ctxExpectation, ctxExpectationAtBegin, ctxExpectationAtEnd int) int {

	// Computing whether the context is clean or not is done by checking whether
	// any batsim or batsched is running. This is intentionally done with a
	// different function (and technique) that the one done within robin.
	ctxCleanAtBegin := batexpe.IsBatsimOrBatschedRunning() == false
	rresult := batexpe.RunRobin(descriptionFile, coverFile, testTimeout)
	ctxCleanAtEnd := batexpe.IsBatsimOrBatschedRunning() == false

	jsonLines, err := batexpe.ParseRobinOutput(rresult.Output)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not parse robin output")
	}

	robintestReturnValue := 0

	// Robin result
	if robinExpectation != EXPECT_NOTHING {
		expectedRobinSuccess := robinExpectation == EXPECT_TRUE
		expectedRobinKilled := robinExpectation == EXPECT_KILLED
		robinSuccess := rresult.ReturnCode == 0
		robinKilled := rresult.Finished == false

		if robinSuccess != expectedRobinSuccess {
			log.WithFields(log.Fields{
				"expected": expectedRobinSuccess,
				"got":      robinSuccess,
			}).Error("Unexpected robin success state")

			robintestReturnValue += UNEXPECTED_ROBIN_SUCCESS_STATE
		}

		if robinKilled != expectedRobinKilled {
			log.WithFields(log.Fields{
				"expected": expectedRobinKilled,
				"got":      robinKilled,
			}).Error("Unexpected robin kill state")

			robintestReturnValue += UNEXPECTED_ROBIN_KILL_STATE
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

			robintestReturnValue += UNEXPECTED_BATSIM_SUCCESS_STATE
		}

		if batsimKilled != expectedBatsimKilled {
			log.WithFields(log.Fields{
				"expected": expectedBatsimKilled,
				"got":      batsimKilled,
			}).Error("Unexpected batsim kill state")

			robintestReturnValue += UNEXPECTED_BATSIM_KILL_STATE
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

			robintestReturnValue += UNEXPECTED_SCHED_SUCCESS_STATE
		}

		if schedPresence != expectedSchedPresence {
			log.WithFields(log.Fields{
				"expected": expectedSchedPresence,
				"got":      schedPresence,
			}).Error("Unexpected sched presence state")

			robintestReturnValue += UNEXPECTED_SCHED_PRESENCE_STATE
		}

		if schedKilled != expectedSchedKilled {
			log.WithFields(log.Fields{
				"expected": expectedSchedKilled,
				"got":      schedKilled,
			}).Error("Unexpected sched kill state")

			robintestReturnValue += UNEXPECTED_SCHED_KILL_STATE
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

			robintestReturnValue += UNEXPECTED_CONTEXT_CLEANLINESS
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

			robintestReturnValue += UNEXPECTED_CONTEXT_CLEANLINESS_AT_BEGIN
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

			robintestReturnValue += UNEXPECTED_CONTEXT_CLEANLINESS_AT_END
		}
	}

	return robintestReturnValue
}
