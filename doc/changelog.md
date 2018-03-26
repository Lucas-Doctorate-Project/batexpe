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
