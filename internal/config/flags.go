// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "flag"

func (c *Config) handleFlagsConfig() {

	flag.Var(&c.Folders, "folders", "Folders or IMAP \"mailboxes\" to check for mail. This value is provided as a comma-separated list.")
	flag.StringVar(&c.Username, "username", "", "The account used to login to the remote mail server. This is often in the form of an email address.")
	flag.StringVar(&c.Password, "password", "", "The remote mail server account password.")
	flag.StringVar(&c.Server, "server", "", "The fully-qualified domain name of the remote mail server.")
	flag.IntVar(&c.Port, "port", 993, "TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections.")
	flag.StringVar(&c.LoggingLevel, "log-level", "info", "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace.")
	flag.BoolVar(&c.EmitBranding, "branding", false, "Toggles emission of branding details with plugin status details. This output is disabled by default.")

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

}
