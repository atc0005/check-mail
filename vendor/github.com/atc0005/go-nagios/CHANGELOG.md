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

## [v0.4.0] - 2020-8-31

### Added

- Add initial "framework workflow"
  - `ExitState` type with `ReturnCheckResults` method
    - used to process and return all applicable check results to Nagios for
      further processing/display
    - supports "branding" callback function to display application name,
      version, or other information as a "trailer" for check results provided
      to Nagios
      - this could be useful for identifying what version of a plugin
        determined the service or host state to be an issue
  - README
    - extend examples to reflect new type/method

### Changed

- GoDoc coverage
  - simple example retained, reader referred to README for further examples

### Fixed

- GitHub Actions Workflow shallow build depth
- `YYYY-MM-DD` changelog version entries

## [v0.3.1] - 2020-08-23

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
  - Add Table of contents

- Linting
  - Local
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
          version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

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

[Unreleased]: https://github.com/atc0005/go-nagios/compare/v0.4.0...HEAD
[v0.4.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.4.0
[v0.3.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.2.0
[v0.1.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.1.0
