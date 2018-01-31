package main

import (
	"context"
	"fmt"
	"github.com/docopt/docopt-go"
	"gitlab.inria.fr/batsim/batexpe"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	EXPECT_NOTHING int = iota
	EXPECT_TRUE
	EXPECT_FALSE
)

type RobinTestResult struct {
	returnCode int
	output     string
}

func main() {
	usage := `Tests one robin execution.

Usage: 
  robintest <description-file>
  			[(--expect-simu-success | --expect-simu-failure)]
  			[(--expect-ctx-clean | --expect-ctx-busy)]
  			[(--expect-batsim-success | --expect-batsim-failure)]
  			[(--expect-sched-success | --expect-sched-failure)]
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
		simuExpectation, ctxExpectation, batsimExpectation,
		schedExpectation, testTimeout)

	os.Exit(testResult)
}

func RobinTest(descriptionFile string,
	simuExpectation, ctxExpectation, batsimExpectation,
	schedExpectation int, testTimeout float64) int {
	return 0
}

func RunRobin(descriptionFile string, testTimeout float64) RobinTestResult {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(testTimeout)*time.Second)
	defer cancel()

	termination := make(chan RobinTestResult)
	executeRobinInnerCtx(ctx, descriptionFile, termination)
	panic("meh")
}

func executeRobinInnerCtx(ctx context.Context, descriptionFile string,
	onexit chan RobinTestResult) {
	cmd := exec.Command("robin")
	cmd.Args = []string{"robin", "--json-logs", descriptionFile}
}
