// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"strings"

	"github.com/atc0005/check-mail/logging"
)

// multiValueFlag is a custom type that satisfies the flag.Value interface in
// order to accept multiple values for some of our flags.
type multiValueFlag []string

// String returns a comma separated string consisting of all slice elements.
func (i *multiValueFlag) String() string {

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if i == nil {
		return ""
	}

	return strings.Join(*i, ", ")
}

// Set is called once by the flag package, in command line order, for each
// flag present.
func (i *multiValueFlag) Set(value string) error {

	// split comma-separated string into multiple folders, toss whitespace
	folders := strings.Split(value, ",")
	for index, folder := range folders {
		folders[index] = strings.TrimSpace(folder)
	}

	// add them to the collection
	*i = append(*i, folders...)
	return nil
}

// Config represents the application configuration as specified via
// command-line flags.
type Config struct {

	// Folders to check for mail. This value is provided a comma-separated
	// list.
	Folders multiValueFlag

	// Username represents the account used to login to the remote mail
	// server. This is often in the form of an email address.
	Username string

	// Password is the remote mail server account password.
	Password string

	// Server is the fully-qualified domain name of the remote mail server.
	Server string

	// Port is the TCP port used to connect to the remote server. This is
	// commonly 993.
	Port int

	// LoggingLevel is the supported logging level for this application.
	LoggingLevel string

	// EmitBranding controls whether "generated by" text is included at the
	// bottom of application output. This output is included in the Nagios
	// dashboard and notifications. This output may not mix well with branding
	// output from other tools such as atc0005/send2teams which also insert
	// their own branding output.
	EmitBranding bool
}

// Version emits application name, version and repo location.
func Version() string {
	return fmt.Sprintf("%s %s (%s)", myAppName, version, myAppURL)
}

// Branding accepts a message and returns a function that concatenates that
// message with version information. This function is intended to be called as
// a final step before application exit after any other output has already
// been emitted.
func Branding(msg string) func() string {
	return func() string {
		return strings.Join([]string{msg, Version()}, "")
	}
}

// Validate verifies all Config struct fields have been provided acceptable
// values.
func (c Config) Validate() error {

	if c.Folders == nil {
		return fmt.Errorf("one or more folders not provided")
	}

	if c.Port < 0 {
		return fmt.Errorf("invalid TCP port number %d", c.Port)
	}

	if c.Username == "" {
		return fmt.Errorf("username not provided")
	}

	if c.Password == "" {
		return fmt.Errorf("password not provided")
	}

	if c.Server == "" {
		return fmt.Errorf("server FQDN not provided")
	}

	requestedLoggingLevel := strings.ToLower(c.LoggingLevel)
	if _, ok := logging.LoggingLevels[requestedLoggingLevel]; !ok {
		return fmt.Errorf("invalid logging level %s", c.LoggingLevel)
	}

	// Optimist
	return nil

}
