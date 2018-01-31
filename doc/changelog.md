# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog][changelog].

Batexpe and all the programs it includes adhere to [Semantic Versioning][semver].  
All the programs share the same version number.

As Batexpe is a library, its public API includes:
- The API of the public functions.

Robin's public API includes:
- The Robin's command-line interface.
- The Robin input file format.

[//]: ==========================================================================
## [Unreleased]
### Added
- New command-line option '--preview-on-error'.  
  If set, robin prints a preview of a process's stdout and stderr on error.  
  This option cannot be set together with `--json-logs`,
  as it directly prints to robin's stdout.

### Fixed
- Regex to find running Batsim processes was bad.

[//]: ==========================================================================
## 0.1.0 - 2018-01-22
- First released version.

[//]: ==========================================================================
[changelog]: http://keepachangelog.com/en/1.0.0/
[semver]: http://semver.org/spec/v2.0.0.html
