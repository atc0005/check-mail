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

## [v0.4.6] - 2021-07-15

### Overview

- Dependency updates
- built using Go 1.16.6

### Changed

- Dependencies
  - `Go`
    - `1.16.5` to `1.16.6`
  - `atc0005/go-nagios`
    - `v0.6.0` to `v0.6.1`
  - `actions/setup-node`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

## [v0.4.5] - 2021-07-06

### Overview

- Add new flag
- Extend initial workaround
- Bug fixes
- built using Go 1.16.5

### Added

- (GH-169, GH-182) Add new `min-tls` flag for specifying minimum supported TLS
  version for use with remote mail server connections

### Changed

- (GH-181) Extend support for user-specified network types (e.g., IPv4-only,
  IPV6-only)

### Fixed

- (GH-184) `net-type` flag keywords comparison is (unintentionally)
  case-sensitive

## [v0.4.4] - 2021-07-04

### Overview

- Bugfix
- built using Go 1.16.5

### Added

- (GH-172, GH-178) Add net-type flag to allow limiting network type to one of
  tcp4 (IPv4-only), tcp6 (IPv6-only) or auto (either).
  - NOTE: This is intended to workaround IPv6 connection failures from an
    IPv4-only Nagios console. See also GH-166 and GH-172.

### Changed

- Documentation
  - minor updates to remove dated information, opting to replace with "see
    also" or general information

## [v0.4.3] - 2021-07-04

### Overview

- Bugfix
- built using Go 1.16.5

### Changed

- Include list of IP Addresses that connection attempts were made against in
  connection failure error message

### Fixed

- String formatting of remote IP Address fails to work for IPv6 addresses

## [v0.4.2] - 2021-07-01

### Overview

- Workarounds
- Dependency updates
- built using Go 1.16.5

### Added

- (GH-167) Create "canary" Dockerfile to track stable Go releases, serve as a
  reminder to generate fresh binaries

### Changed

- (GH-166) Attempt to connect to each server IP before aborting connection
  attempt
  - NOTE: This is intended to workaround IPv6 connection failures from an
    IPv4-only Nagios console. See also GH-169 and GH-172.
- (GH-166, GH-170) Explicitly require TLS minimum version of `v1.2` over
  default of `v1.0`
  - NOTE: Per GH-169, this is intended to be user-configurable in the future

- Dependencies
  - built using Go 1.16.5
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `actions/setup-node`
    - `v2.1.5` to `v2.2.0`
  - `rs/zerolog`
    - `v1.22.0` to `v1.23.0`

## [v0.4.1] - 2021-05-15

### Overview

- Bug fix
- Dependency updates
- built using Go 1.16.4

### Changed

- Dependencies
  - built using Go 1.16.4
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `emersion/go-imap`
    - `v1.0.6` to `v1.1.0`
  - `rs/zerolog`
    - `v1.21.0` to `v1.22.0`

### Fixed

- list-emails | Emails containing hash/pound/number sign characters are not
  escaped
  - reported by @dmd0032 (along with suggested HTML entity replacement)

## [v0.4.0] - 2021-04-15

### Overview

- Bug fixes
- Misc tweaks
- built using Go 1.16.3

### Added

- Add support for escaping Textile formatting characters

### Changed

- Dependencies
  - built using Go 1.16.3
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `actions/setup-node`
    - `v2.1.4` to `v2.1.5`
  - `rs/zerolog`
    - `v1.20.0` to `v1.21.0`

### Fixed

- list-emails | Emails containing pipe characters are not properly escaped
- Fix doc comment string typos

## [v0.3.3] - 2021-02-15

### Overview

- Bug fixes
- Misc tweaks
- built using Go 1.15.8

### Changed

- Swap out GoDoc badge for pkg.go.dev badge
- Dependencies
  - built using Go 1.15.8
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`
  - `atc0005/go-nagios`
    - `v0.5.2` to `v0.6.0`

## [v0.3.2] - 2020-12-09

### Changed

- Extend Unicode string tests
  - new entries
  - refactored to store all test strings separately for reuse

- Minor template tweaks
  - New `Reported on` column in report
  - Tweak description of placeholder character in an attempt to better
    communicate purpose

- Dependencies
  - built using Go 1.15.6
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)

## [v0.3.1] - 2020-11-30

### Overview

- Bugfixes

- Built using Go 1.15.5
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Fixed

- Validation of log file and report output directories fail for Nagios plugin
- Sub-mailboxes beneath "Inbox" require uppercase "INBOX" name as prefix in
  order to reference them

## [v0.3.0] - 2020-11-26

### Overview

- Add Go version of [`list-emails`](https://github.com/atc0005/list-emails)
  Python script

- Built using Go 1.15.5
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

- Significant refactoring

### Added

- Port [`list-emails`](https://github.com/atc0005/list-emails) Python script

### Changed

- Internal: significant refactoring to share code with `list-emails` CLI
  application

### Fixed

- Fix breadcrumb URLs
- Fix ServiceOutput templating for multiple accounts

## [v0.2.9] - 2020-11-23

### Overview

- Built using Go 1.15.5
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

- Minor refactoring
- Minor code polish

### Changed

- Update Nagios plugin to support checking multiple accounts
  - this isn't currently exposed (via command-line flags or otherwise)
  - this was added to support upcoming inclusion of the `list-emails`
    application which will process accounts in bulk

### Fixed

- Minor "happy path" fix
- Update logging for deferred logout step

## [v0.2.8] - 2020-11-17

### Overview

- Built using Go 1.15.5
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

- Minor changes overall, mostly setup work for further development.

### Added

- Add version flag

### Changed

- Internal refactoring of config initialization, logging and flags

## [v0.2.7] - 2020-11-16

### Changed

- Dependencies
  - built using Go 1.15.5
    - **Statically linked**
    - Windows
      - x86
      - x64
    - Linux
      - x86
      - x64
  - `atc0005/go-nagios`
    - `v0.5.1` to `v0.5.2`
  - `emersion/go-imap`
    - `v1.0.5` to `v1.0.6`
  - `actions/checkout`
    - `v2.3.3` to `v2.3.4`

### Fixed

- Add additional context to `v0.2.6` release notes

## [v0.2.6] - 2020-10-06

Follow-up binary release built using Go 1.15.2, cgo explicitly disabled. See
GH-101 for more details.

### Changed

- Add `-trimpath` build flag
  - intended to help prune verbose/unnecessary details from output

### Fixed

- Makefile build options generate static binaries which *potentially* bundle
  glibc (LGPL license)
  - I have been unable to confirm this, so attempting to "play it safe"
  - the goal is to avoid attaching the LGPL license to this project until I've
    properly researched and understand the ramifications of doing so for this
    project, for future forks, etc.

## [v0.2.5] - 2020-10-03

### Added

- First (limited) binary release
  - built using Go 1.15.2
    - see GH-94 and GH-95
  - windows
    - x86
    - x64
  - linux
    - x86
    - x64

### Changed

- Dependencies
  - built using Go 1.15.2
    - see GH-94 and GH-95
  - `atc0005/go-nagios`
    - `v0.4.0` to `v0.5.1`
  - `rs/zerolog`
    - `v1.19.0` to `v1.20.0`
  - `actions/checkout`
    - `v2.3.2` to `v2.3.3`
  - `actions/setup-node`
    - `v2.1.1` to `v2.1.2`

- Move `config`, `mailboxes` code into subpackages

- `ReturnNagiosResults` deferred first, allowed to run last (as intended) to
  handle setting final exit code
  - this was changed to match new functionality added in `atc0005/go-nagios`
    dependency

### Fixed

- Makefile build options do not generate static binaries
- Typo in CHANGELOG release date format
- Makefile generates checksums with qualified path
- Failure to initialize application configuration at startup; "Error
  validating configuration: one or more folders not provided" message printed

## [v0.2.4] - 2020-08-31

### Changes

- Dependencies
  - upgrade `atc0005/go-nagios`
    - `v0.3.1` to `v0.4.0`

- Integration changes to replace custom check results handling with
  functionality provided by `v0.4.0` of the `atc0005/go-nagios` package

## [v0.2.3] - 2020-08-23

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable"
      - "stable"
      - "unstable"
    - Makefile is *not* used in this workflow
    - staticcheck linting using latest stable version provided by the
      `atc0005/go-ci` containers

  - Separate Makefile-based linting and building workflow
    - intended to help ensure that local Makefile-based builds that are
      referenced in project README files continue to work as advertised until
      a better local tool can be discovered/explored further
    - use `golang:latest` container to allow for Makefile-based linting
      tooling installation testing since the `atc0005/go-ci` project provides
      containers with those tools already pre-installed
      - linting tasks use container-provided `golangci-lint` config file
        *except* for the Makefile-driven linting task which continues to use
        the repo-provided copy of the `golangci-lint` configuration file

  - Add Quick Validation workflow
    - run on every push, everything else on pull request updates
    - linting via `golangci-lint` only
    - testing
    - no builds

- Add new README badges for additional CI workflows
  - each badge also links to the associated workflow results

### Changed

- Disable `golangci-lint` default exclusions

- dependencies
  - `go.mod` Go version
    - updated from `1.13` to `1.14`
  - `actions/setup-go`
    - updated from `v2.1.0` to `v2.1.2`
      - since replaced with Docker containers
  - `actions/setup-node`
    - updated from `v2.1.0` to `v2.1.1`
  - `actions/checkout`
    - updated from `v2.3.1` to `v2.3.2`

- README
  - Link badges to applicable GitHub Actions workflows results

- Linting
  - Local
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
          version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

## [v0.2.2] - 2020-07-05

### Changed

- Dependencies
  - `actions/setup-go`
    - `v2.0.3` to `v2.1.0`
  - `actions/setup-node`
    - `v2.0.0` to `v2.1.0`
  - `atc0005/go-nagios`
    - `v0.2.0` to `v0.3.0`

- Replace hard-coded status strings with const refs
  - New in `v0.3.0` of `atc0005/go-nagios`

- Fix exit code references
  - These have changed as of `v0.3.0` of `atc0005/go-nagios`

- Minor README tweaks
  - Fix path to generated binaries
  - Markdown VSCode extension (`yzhang.markdown-all-in-one`) auto-fixes to ToC

## [v0.2.1] - 2020-06-23

### Added

- Enable Dependabot updates
  - Go Modules
  - GitHub Actions

### Changed

- Dependencies
  - `rs/zerolog`
    - `v1.18.0` to `v1.19.0`
  - `emersion/go-imap`
    - `v1.0.4` to `v1.0.5`
  - `actions/setup-go`
    - `v1` to `v2.0.3`
  - `actions/checkout`
    - `v1` to `v2.3.1`
  - `actions/setup-node`
    - `v1` to `v2.0.0`
  - `golangci-lint`
    - `1.25.0` to `v1.27.0`

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

[Unreleased]: https://github.com/atc0005/check-mail/compare/v0.4.6...HEAD
[v0.4.6]: https://github.com/atc0005/check-mail/releases/tag/v0.4.6
[v0.4.5]: https://github.com/atc0005/check-mail/releases/tag/v0.4.5
[v0.4.4]: https://github.com/atc0005/check-mail/releases/tag/v0.4.4
[v0.4.3]: https://github.com/atc0005/check-mail/releases/tag/v0.4.3
[v0.4.2]: https://github.com/atc0005/check-mail/releases/tag/v0.4.2
[v0.4.1]: https://github.com/atc0005/check-mail/releases/tag/v0.4.1
[v0.4.0]: https://github.com/atc0005/check-mail/releases/tag/v0.4.0
[v0.3.3]: https://github.com/atc0005/check-mail/releases/tag/v0.3.3
[v0.3.2]: https://github.com/atc0005/check-mail/releases/tag/v0.3.2
[v0.3.1]: https://github.com/atc0005/check-mail/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/check-mail/releases/tag/v0.3.0
[v0.2.9]: https://github.com/atc0005/check-mail/releases/tag/v0.2.9
[v0.2.8]: https://github.com/atc0005/check-mail/releases/tag/v0.2.8
[v0.2.7]: https://github.com/atc0005/check-mail/releases/tag/v0.2.7
[v0.2.6]: https://github.com/atc0005/check-mail/releases/tag/v0.2.6
[v0.2.5]: https://github.com/atc0005/check-mail/releases/tag/v0.2.5
[v0.2.4]: https://github.com/atc0005/check-mail/releases/tag/v0.2.4
[v0.2.3]: https://github.com/atc0005/check-mail/releases/tag/v0.2.3
[v0.2.2]: https://github.com/atc0005/check-mail/releases/tag/v0.2.2
[v0.2.1]: https://github.com/atc0005/check-mail/releases/tag/v0.2.1
[v0.2.0]: https://github.com/atc0005/check-mail/releases/tag/v0.2.0
[v0.1.2]: https://github.com/atc0005/check-mail/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-mail/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-mail/releases/tag/v0.1.0
