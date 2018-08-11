package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"gitlab.inria.fr/batsim/batexpe"
	"io/ioutil"
	"os"
	"strconv"
)

func setupLogging(arguments map[string]interface{}) (previewOnError bool) {
	log.SetOutput(os.Stdout)

	if arguments["--json-logs"] == true {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		customFormatter := new(log.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
		customFormatter.FullTimestamp = true
		customFormatter.QuoteEmptyFields = true
		log.SetFormatter(customFormatter)
	}

	if arguments["--debug"] == true {
		log.SetLevel(log.DebugLevel)
	} else if arguments["--quiet"] == true {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	previewOnError = false
	if arguments["--preview-on-error"] == true {
		previewOnError = true
	}

	return previewOnError
}

func ExperimentFromArgs(arguments map[string]interface{}) batexpe.Experiment {
	var exp batexpe.Experiment
	var err error

	// Default values
	exp.Batcmd = "batcmd-unset"
	exp.OutputDir = "output-dir-unset"
	exp.Schedcmd = "schedcmd-unset"
	exp.SimulationTimeout = 604800
	exp.ReadyTimeout = 10
	exp.SuccessTimeout = 3600
	exp.FailureTimeout = 5

	if arguments["--batcmd"] != nil {
		exp.Batcmd = arguments["--batcmd"].(string)
	}

	if arguments["--output-dir"] != nil {
		exp.OutputDir = arguments["--output-dir"].(string)
	}

	if arguments["--schedcmd"] != nil {
		exp.Schedcmd = arguments["--schedcmd"].(string)
	}

	if arguments["--simulation-timeout"] != nil {
		exp.SimulationTimeout, err = strconv.ParseFloat(
			arguments["--simulation-timeout"].(string), 64)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"--simulation-timeout": arguments["--simulation-timeout"].(string),
			}).Fatal("Invalid simulation timeout")
		}
	}

	if arguments["--ready-timeout"] != nil {
		exp.ReadyTimeout, err = strconv.ParseFloat(
			arguments["--ready-timeout"].(string), 64)
		if err != nil {
			log.WithFields(log.Fields{
				"err":             err,
				"--ready-timeout": arguments["--ready-timeout"].(string),
			}).Fatal("Invalid ready timeout")
		}
	}

	if arguments["--success-timeout"] != nil {
		exp.SuccessTimeout, err = strconv.ParseFloat(
			arguments["--success-timeout"].(string), 64)
		if err != nil {
			log.WithFields(log.Fields{
				"err":               err,
				"--success-timeout": arguments["--success-timeout"].(string),
			}).Fatal("Invalid success timeout")
		}
	}

	if arguments["--failure-timeout"] != nil {
		exp.FailureTimeout, err = strconv.ParseFloat(
			arguments["--failure-timeout"].(string), 64)
		if err != nil {
			log.WithFields(log.Fields{
				"err":               err,
				"--failure-timeout": arguments["--failure-timeout"].(string),
			}).Fatal("Invalid failure timeout")
		}
	}

	log.WithFields(log.Fields{
		"args": arguments,
		"expe": exp,
	}).Debug("arguments -> expe")

	return exp
}

func generateDescription(arguments map[string]interface{}) {
	exp := ExperimentFromArgs(arguments)
	yam := batexpe.ToYaml(exp)

	fil := arguments["<description-file>"].(string)

	err := ioutil.WriteFile(fil, []byte(yam), 0644)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"filename": fil,
		}).Fatal("Cannot write file")
	}
}

func main() {
	os.Exit(mainReturnWithCode())
}

func mainReturnWithCode() int {
	usage := `Robin manages the execution of one Batsim simulation.

Usage:
  robin --output-dir=<dir>
        --batcmd=<batsim-command>
        [--schedcmd=<scheduler-command>]
        [--simulation-timeout=<time>]
        [--ready-timeout=<time>]
        [--success-timeout=<time>]
        [--failure-timeout=<time>]
        [(--verbose | --quiet | --debug)] [(--json-logs | --preview-on-error)]
  robin <description-file>
        [(--verbose | --quiet | --debug)] [(--json-logs | --preview-on-error)]
  robin generate <description-file>
        [--output-dir=<dir>]
        [--batcmd=<batsim-command>]
        [--schedcmd=<scheduler-command>]
        [--simulation-timeout=<time>]
        [--ready-timeout=<time>]
        [--success-timeout=<time>]
        [--failure-timeout=<time>]
        [(--verbose | --quiet | --debug)] [(--json-logs | --preview-on-error)]
  robin -h | --help
  robin --version


Examples:
  robin --output-dir=/tmp \
        --batcmd="batsim -p platform.xml -w workload.json" \
        --schedcmd="batsched"
  robin --output-dir=/tmp \
        --batcmd="batsim -p platform.xml -w workload.json --batexec"
  robin input_description_file.yaml
  robin generate output_description_file.yaml


Timeout options:
  --simulation-timeout=<time>   Simulation timeout in seconds.
                                If this time is exceeded, the simulation is
                                stopped. Default value is one week.
                                [default: 604800]

  --ready-timeout=<time>        Ready timeout in seconds.
                                If the context is still invalid after this
                                time, the script is stopped.
                                This includes socket already in use and
                                conflicting Batsim instances.
                                [default: 10]

  --success-timeout=<time>      The timeout for the second process to complete
                                once the first process has finished
                                successfully (returned 0).
                                [default: 3600]

  --failure-timeout=<time>      The timeout for the second process to complete
                                once the first process has finished
                                unsuccessfully.
                                [default: 5]

Verbosity options:
  --quiet                       Only print critical information.
  --verbose                     Print information. Default verbosity mode.
  --debug                       Print debug information.
  --json-logs                   Print information in JSON.
  --preview-on-error            Preview stdout and stderr of failed processes.`

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

	arguments, err := parser.ParseArgs(usage, os.Args[1:], batexpe.Version())
	if ret != -1 {
		return ret
	}
	if err != nil {
		return 2
	}

	previewOnError := setupLogging(arguments)

	log.WithFields(log.Fields{
		"args": arguments,
	}).Debug("Arguments parsed")

	// Generate mode?
	if arguments["generate"] == true {
		generateDescription(arguments)
	}

	// Execution mode.
	// Read what should be executed
	var exp batexpe.Experiment
	if arguments["<description-file>"] != nil {
		fil := arguments["<description-file>"].(string)
		byt, err := ioutil.ReadFile(fil)
		if err != nil {
			log.WithFields(log.Fields{
				"err":      err,
				"filename": fil,
			}).Fatal("Cannot open description file")
		}

		exp = batexpe.FromYaml(string(byt))
	} else {
		exp = ExperimentFromArgs(arguments)
	}

	log.WithFields(log.Fields{
		"batsim command":     exp.Batcmd,
		"output directory":   exp.OutputDir,
		"scheduler command":  exp.Schedcmd,
		"simulation timeout": exp.SimulationTimeout,
		"ready timeout":      exp.ReadyTimeout,
		"success timeout":    exp.SuccessTimeout,
		"failure timeout":    exp.FailureTimeout,
	}).Debug("Instance description read")

	ret = batexpe.ExecuteOne(exp, previewOnError)
	return ret
}
