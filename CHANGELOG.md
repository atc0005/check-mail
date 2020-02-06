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

[Unreleased]: https://github.com/atc0005/check-mail/compare/v0.1.2...HEAD
[v0.1.2]: https://github.com/atc0005/check-mail/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-mail/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-mail/releases/tag/v0.1.0
