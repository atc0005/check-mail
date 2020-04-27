#!/bin/bash

# Copyright 2020 Adam Chalkley
#
# https://github.com/atc0005/check-mail
#
# Licensed under the MIT License. See LICENSE file in the project root for
# full license information.


# Purpose: Helper script for installing linting tools used by this project

export PATH=${PATH}:$(go env GOPATH)/bin

# https://github.com/golangci/golangci-lint#install
# https://github.com/golangci/golangci-lint/releases/latest
GOLANGCI_LINT_VERSION="v1.25.0"

# Temporarily disable module-aware mode so that we can install linting tools
# without modifying this project's go.mod and go.sum files
export GO111MODULE="off"
go get -u golang.org/x/lint/golint
go get -u honnef.co/go/tools/cmd/staticcheck

echo Installing golangci-lint ${GOLANGCI_LINT_VERSION} per official binary installation docs ...
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}
golangci-lint --version

# Reset GO111MODULE back to the default
export GO111MODULE=""
