# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog][changelog].

Batexpe and all the programs it includes adhere to
[Semantic Versioning][semver].  
All the programs share the same version number.

As Batexpe is a library, its public API includes:
- The API of the public functions.

Robin's public API includes:
- The Robin's command-line interface.
- The Robin input file format.

Robintest's public API includes:
- The Robintest's command-line interface.

[//]: ==========================================================================
## [Unreleased]
### Changed
- Batsim commands are no longer directly parsed by batexpe.
  This is now delegated to Batsim by calling it with
  the `--dump-execution-context` command-line flag.
- Robin now takes Redis into account to detect whether two batsim instances
  are conflicting or not. Previously, only the port used between Batsim and
  the scheduler was checked.

[//]: ==========================================================================
## [1.1.0] - 2018-08-11
### Added
- New `--no-preview-on-error` option,
  that disables the preview of the logs of failed processes.

### Changed
- Improved the preview of the logs of failed processes:
  - This is now written on stderr (was on stdout).
  - This is now compatible with the `--json-logs` option, as stdout can keep
    its JSON structure regardless of stderr previews.
  - This is now enabled by default
    (use `--no-preview-on-error` to keep old behaviour).

### Deprecated
- The `--preview-on-error` option is now deprecated.

[//]: ==========================================================================
## [1.0.0] - 2018-04-12
### Added
- The version given by `--version` can now be set from `git describe`.  
  For build convenience's sake, a `Makefile` is provided at the
  repository's root.

### Changed
- Most features are now tested and seem to work.  
  All code is covered (maybe not fully on the CI yet).
- All calls to `log.Fatal` or `os.Exit` have been removed.
  The following functions now also return an error:
  - `FromYaml`, `IsBatsimOrBatschedRunning`, `PrepareDirs`, `ToYaml`
- robintest now consider that robin's execution context was clean during
  execution unless it encounters a log line stating the opposite.  
  Previously, robintest expected the execution context to be busy by default
  (unless encountering a log line stating the opposite).

### Fixed
- Setgpid was not set on some user-given commands (batsim command when batsim
  was launched without scheduler, and the check script).
  This resulted in Kill not working as expected (only the subprocess was
  killed, not the subprocess and all its children).

[//]: ==========================================================================
## [0.3.0] - 2018-04-08
### Added
- More tests, coverage reports in CI.
- robintest can now execute a check script with the ``--result-check-script``
  option. Such script is called if all expectations have been met. The
  script is called with batsim's export prefix as first argument.
- robintest now supports a ``--debug`` logging option.
- robintest can now be built with coverage support.
- robintest can now call ``robin.test`` (robin with coverage support) via
  the ``--cover`` option.
- More batexpe functions and types are now public:
  - Types: ``CmdFinishedMsg``
  - Functions: ``ExecuteTimeout``

### Changed
- Most robin functions should now return an error.
- ``RobinResult`` struct now has a ``Succeeded`` boolean instead of a
  ``ReturnCode`` integer.

### Fixed
- robin tried to execute the instance even with the ``generate`` subcommand.  
  robin should now return after generating the description file.
- robintest return value could be 0 while an expection was not met.  
  robintest should now return 1 when expectations are not met.
- robintest did not retrieve robin's return code correctly.
- batexpe code to determine whether conflicting batsim instances are running
  did not work: Some batsim instances were not found.  
  The regexp has been improved and they seem to be detected now.

[//]: ==========================================================================
## [0.2.0] - 2018-03-27
### Added
- Robin: New command-line option '--preview-on-error'.  
  If set, robin prints a preview of a process's stdout and stderr on error.  
  This option cannot be set together with `--json-logs`,
  as it directly prints to robin's stdout.
- New ``robintest`` program, meant to wrap robin calls, parse their output
  and check that expected behaviors happened.
- Robin can now be built with coverage support.
  As compiling and running robin with coverage can be tricky, please refer to
  [batexpe's CI script](../.gitlab-ci.yml) for more information.
- [Robin's CI](https://gitlab.inria.fr/batsim/batexpe/pipelines) has been set
  up. Robin it heavily tested for simple cases and should now work for them.
- Batexpe: New function ``KillProcess``, that kills a process and all its
  children (sending SIGTERM to them).
- Batexpe: New function ``PreviewFile``, that reads a file and display a
  preview of its content, showing the whole file if short enough or only its
  first and last lines.
- Batexpe: New function ``IsBatsimOrBatschedRunning``, that calls ``ps`` and
  parse its output to determine whether any batsim or batsched is running.
- Batexpe: New types and functions to ease parsing robin's output
  (in [parserobin.go](../parserobin.go):
  - Types: RobinResult
  - Functions: ``RunRobin``, ``ParseRobinOutput``, ``WasBatsimSuccessful``,
  ``WasSchedSuccessful``, ``WasContextClean``.

### Changed
- Batexpe's ``ParseBatsimCommand`` now also returns an error.
- Batexpe's ``CreateDirIfNeeded`` now returns an error.
- Batexpe's ``PortFromBatSock`` now also returns an error.
  It now also expects a batsim command as input parameter.

### Fixed
- Regex to find running Batsim processes was bad.
- Typing ctrl+C too fast or setting very low timeouts caused segmentation fault
  when killing processes. This should now be fixed.
- Process cleanup robustness has been improved.
  It should now work in most simple cases, as seen in
  [Robin's CI](https://gitlab.inria.fr/batsim/batexpe/pipelines).

[//]: ==========================================================================
## 0.1.0 - 2018-01-22
- First released version.

[//]: ==========================================================================
[changelog]: http://keepachangelog.com/en/1.0.0/
[semver]: http://semver.org/spec/v2.0.0.html

[Unreleased]: https://framagit.org/batsim/batexpe/compare/v1.1.0...master
[1.1.0]: https://framagit.org/batsim/batexpe/compare/v1.0.0...v1.1.0
[1.0.0]: https://framagit.org/batsim/batexpe/compare/v0.3.0...v1.0.0
[0.3.0]: https://framagit.org/batsim/batexpe/compare/v0.2.0...v0.3.0
[0.2.0]: https://framagit.org/batsim/batexpe/compare/v0.1.0...v0.2.0
