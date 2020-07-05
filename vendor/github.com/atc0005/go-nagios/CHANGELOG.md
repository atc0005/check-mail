# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/go-nagios/issues) for any
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

## [v0.3.0] - 2020-07-05

### Added

- Add State "labels" to provide an alternative to using literal state strings

- Add GitHub Actions Workflow, Makefile for builds
  - Lint codebase
  - Build codebase

- Enable Dependabot updates
  - GitHub Actions
  - Go Modules

### Changed

- BREAKING: Rename existing exit code constants to explicitly note that they
  are exit codes
  - the thinking was that since we have text "labels" for state, it would be
    good to help clarify the difference between the new constants and the
    existing exit code constants

- Minor tweaks to README to reference changes, wording

- Update dependencies
  - `actions/checkout`
    - `v1` to `v2.3.1`
  - `actions/setup-go`
    - `v2.0.3` to `v2.1.0`
  - `actions/setup-node`
    - `v1` to `v2.1.0`

## [v0.2.0] - 2020-02-02

### Added

- Add Nagios `State` constants

### Removed

- Nagios `State` map

## [v0.1.0] - 2020-01-20

### Added

Initial package state

- Nagios state map

[Unreleased]: https://github.com/atc0005/go-nagios/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.2.0
[v0.1.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.1.0
