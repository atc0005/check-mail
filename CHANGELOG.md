# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/check-mail/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.2.0] - 2020-04-28

### Changed

- Use common `cmd` subdirectory structure in order to more easily support
  multiple binaries

- Vendor dependencies

- Update Makefile and helper shell scripts
  - include `-mod=vendor` build flag for applicable `go` commands to reflect
    Go 1.13 vendoring
    - this includes specifying `-mod=vendor` even for `go list` commands,
        which unless specified results in dependencies being downloaded, even
        when they're already provided in a local, top-level `vendor` directory

- Update GitHub Actions Workflows
  - Disable running `go get` after checking out code
  - Exclude `vendor` folder from ...
    - Markdown linting checks
    - tests
    - basic build
  - Echo Go version used for CI runs
  - Update Go versions used
    - Remove Go 1.12 (no longer supported)
    - Add Go 1.14 (recent release)

- Dependencies
  - Update rs/zerolog to v1.18.0
  - Update emersion/go-imap to v1.0.4

- Linting
  - Move golangci-lint settings to external file
  - Add scopelint golangci-lint linter
  - Use golangci-lint binary instead of building src
  - Replace external shell scripts by incorporating applicable commands
    directly into the Makefile
  - Disable gofmt, golint external commands, rely on golangci-lint for that
    linting functionality

- Documentation
  - Update README to reflect recent updates to build process/layout

## [v0.1.2] - 2020-02-06

### Fixed

- Update status output to reflect the same format used in the original Python
  2 plugin.
  - For reasons I've yet to spend sufficient time to figure out, the
    double-quoting used for elements of the "folders" list is lost when sent
    by Teams or email notifications. It is easier to go ahead and just revert
    the format for now so it is consistent in each format (console, Teams or
    email).

## [v0.1.1] - 2020-02-06

### Fixed

- Branding output (app name, version) was only shown for application error
  conditions. This has been adjusted so that it is intentionally not shown by
  default for any condition, but can be toggled on via a new `-branding` flag.
- README example

## [v0.1.0] - 2020-02-04

Initial release!

This release provides an early release version of a Nagios plugin used to
monitor IMAP mailboxes for content. This plugin (or its predecessor as of this
writing) is used to monitor email accounts scraped by our ticketing system.

Future releases of this project are expected to shift directory structure and
content in order to accommodate additional Nagios plugins and tools used to
monitor mail-related resources.

### Added

- Monitor one or many mailboxes
- Optional, leveled logging using `rs/zerolog` package
  - JSON-format output
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`.
- TLS/SSL IMAP4 connectivity via `emerson/go-imap` package
- Go modules (vs classic `GOPATH` setup)

[Unreleased]: https://github.com/atc0005/check-mail/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/atc0005/check-mail/releases/tag/v0.2.0
[v0.1.2]: https://github.com/atc0005/check-mail/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-mail/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-mail/releases/tag/v0.1.0
