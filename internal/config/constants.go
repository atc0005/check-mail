// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

// help text for our CLI flags, maintained in one common block
const (
	foldersFlagHelp      string = "Folders or IMAP \"mailboxes\" to check for mail. This value is provided as a comma-separated list."
	usernameFlagHelp     string = "The account used to login to the remote mail server. This is often in the form of an email address."
	passwordFlagHelp     string = "The remote mail server account password."
	serverFlagHelp       string = "The fully-qualified domain name of the remote mail server."
	portFlagHelp         string = "TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections."
	loggingLevelFlagHelp string = "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace."
	emitBrandingFlagHelp string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	versionFlagHelp      string = "Whether to display application version and then immediately exit application."
)

// Default flag settings if not overridden by user input
const (
	defaultLoggingLevel          string = "info"
	defaultEmitBranding          bool   = false
	defaultPort                  int    = 993
	defaultServer                string = ""
	defaultPassword              string = ""
	defaultUsername              string = ""
	defaultDisplayVersionAndExit bool   = false
)
