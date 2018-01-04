package batexpe

import (
	"github.com/anmitsu/go-shlex"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

type BatsimArgs struct {
	Socket       string
	ExportPrefix string
	BatexecMode  bool
}

func ParseBatsimCommand(batcmd string) BatsimArgs {
	batsimDocopt := `
A tool to simulate (via SimGrid) the behaviour of scheduling algorithms.

Usage:
  batsim -p <platform_file> [-w <workload_file>...]
                            [-W <workflow_file>...]
                            [--WS (<cut_workflow_file> <start_time>)...]
                            [options]
  batsim --help
  batsim --version
  batsim --simgrid-version
  batsim --unittest

Input options:
  -p --platform <platform_file>     The SimGrid platform to simulate.
  -w --workload <workload_file>     The workload JSON files to simulate.
  -W --workflow <workflow_file>     The workflow XML files to simulate.
  --WS --workflow-start (<cut_workflow_file> <start_time>)... The workflow XML
                                    files to simulate, with the time at which
                                    they should be started.

Most common options:
  -m, --master-host <name>          The name of the host in <platform_file>
                                    which will be used as the RJMS management
                                    host (thus NOT used to compute jobs)
                                    [default: master_host].
  -E --energy                       Enables the SimGrid energy plugin and
                                    outputs energy-related files.

Execution context options:
  --config-file <cfg_file>          Configuration file name (optional). [default: None]
  -s, --socket-endpoint <endpoint>  The Decision process socket endpoint
                                    Decision process [default: tcp://localhost:28000].
  --redis-hostname <redis_host>     The Redis server hostname. Read from config file by default.
                                    [default: None]
  --redis-port <redis_port>         The Redis server port. Read from config file by default.
                                    [default: -1]
  --redis-prefix <prefix>           The Redis prefix. Read from config file by default.
                                    [default: None]

Output options:
  -e, --export <prefix>             The export filename prefix used to generate
                                    simulation output [default: out].
  --enable-sg-process-tracing       Enables SimGrid process tracing
  --disable-schedule-tracing        Disables the Pajé schedule outputting.
  --disable-machine-state-tracing   Disables the machine state outputting.


Platform size limit options:
  --mmax <nb>                       Limits the number of machines to <nb>.
                                    0 means no limit [default: 0].
  --mmax-workload                   If set, limits the number of machines to
                                    the 'nb_res' field of the input workloads.
                                    If several workloads are used, the maximum
                                    of these fields is kept.
Verbosity options:
  -v, --verbosity <verbosity_level> Sets the Batsim verbosity level. Available
                                    values: quiet, network-only, information,
                                    debug [default: information].
  -q, --quiet                       Shortcut for --verbosity quiet

Workflow options:
  --workflow-jobs-limit <job_limit> Limits the number of possible concurrent
                                    jobs for workflows. 0 means no limit
                                    [default: 0].
  --ignore-beyond-last-workflow     Ignores workload jobs that occur after all
                                    workflows have completed.

Other options:
  --allow-time-sharing              Allows time sharing: One resource may
                                    compute several jobs at the same time.
  --batexec                         If set, the jobs in the workloads are
                                    computed one by one, one after the other,
                                    without scheduler nor Redis.
  --pfs-host <pfs_host>             The name of the host, in <platform_file>,
                                    which will be the parallel filesystem target
                                    as data sink/source for the large-capacity
                                    storage tier [default: pfs_host].
  --hpst-host <hpst_host>           The name of the host, in <platform_file>,
                                    which will be the parallel filesystem target
                                    as data sink/source for the high-performance
                                    storage tier [default: hpst_host].
  -h --help                         Shows this help.
`

	splitCmd, err := shlex.Split(batcmd, true)
	if err != nil {
		log.WithFields(log.Fields{
			"err":            err,
			"batsim command": batcmd,
		}).Fatal("Cannot split Batsim command")
	}

	log.WithFields(log.Fields{
		"batsim command": batcmd,
	}).Debug("Batcmd -> shlex -> split")

	if len(splitCmd) <= 1 {
		log.WithFields(log.Fields{
			"batsim command": batcmd,
		}).Fatal("Batsim command should have at least 2 parts")
	}

	arguments, err := docopt.Parse(batsimDocopt, splitCmd[1:], true, "?",
		false, false)
	if err != nil {
		log.WithFields(log.Fields{
			"err":            err,
			"batsim command": batcmd,
		}).Fatal("Cannot parse Batsim command")
	}

	var batargs BatsimArgs
	batargs.Socket = "tcp://localhost:28000"
	batargs.ExportPrefix = "out"
	batargs.BatexecMode = false

	if arguments["--socket-endpoint"] != nil {
		batargs.Socket = arguments["--socket-endpoint"].(string)
	}

	if arguments["--export"] != nil {
		batargs.ExportPrefix = arguments["--export"].(string)
	}

	if arguments["--batexec"] != nil {
		batargs.BatexecMode = arguments["--batexec"].(bool)
	}

	log.WithFields(log.Fields{
		"socket":        batargs.Socket,
		"export prefix": batargs.ExportPrefix,
		"batexec?":      batargs.BatexecMode,
	}).Debug("Parsed Batsim command")

	return batargs
}
