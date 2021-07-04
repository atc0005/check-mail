// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "os"

// help text for our CLI flags, maintained in one common block
const (
	foldersFlagHelp             string = "Folders or IMAP \"mailboxes\" to check for mail. This value is provided as a comma-separated list."
	usernameFlagHelp            string = "The account used to login to the remote mail server. This is often in the form of an email address."
	passwordFlagHelp            string = "The remote mail server account password."
	serverFlagHelp              string = "The fully-qualified domain name of the remote mail server."
	portFlagHelp                string = "TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections."
	networkTypeFlagHelp         string = "Limits network connections to remote mail servers to one of tcp4 (IPv4-only), tcp6 (IPv6-only) or auto (either)."
	loggingLevelFlagHelp        string = "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace."
	emitBrandingFlagHelp        string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	versionFlagHelp             string = "Whether to display application version and then immediately exit application."
	iniConfigFileFlagHelp       string = "Full path to the INI-formatted configuration file used by this application. See contrib/list-emails/accounts.example.ini for a starter template. Rename to accounts.ini, update with applicable information and place in a directory of your choice. If this file is found in your current working directory you need not use this flag."
	reportFileOutputDirFlagHelp string = "Full path to the directory where email summary report files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist."
	logFileOutputDirFlagHelp    string = "Full path to the directory where log files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist."
)

// Default flag settings if not overridden by user input
const (
	defaultLoggingLevel          string = "info"
	defaultEmitBranding          bool   = false
	defaultPort                  int    = 993
	defaultServer                string = ""
	defaultPassword              string = ""
	defaultUsername              string = ""
	defaultNetworkType           string = netTypeTCPAuto
	defaultDisplayVersionAndExit bool   = false

	// By default these directories are created/used in the user's current
	// working directory. The workflow for the older, Python-based list-emails
	// script was to create these directories in the same directory as the
	// script itself. If the application happens to be in the current working
	// directory, the resulting behavior will be the same.
	defaultReportFileOutputDir string = "output"
	defaultLogFileOutputDir    string = "log"

	// defaultINIConfigFileName is the "bare" or "non-qualified" configuration
	// filename. If not specified explicitly via CLI flag, this file will be
	// looked for alongside the app's location.
	defaultINIConfigFileName string = "accounts.ini"

	// defaultTOMLConfigFileName is the "bare" or "non-qualified"
	// configuration filename. If not specified explicitly via CLI flag, this
	// file will be looked for alongside the app's location first. If not
	// found, then the config path for the user account running this app will
	// be checked, e.g., `$HOME/.check-mail/config.toml`. This file is
	// automatically created alongside this application from an existing
	// `accounts.ini` file, if found. If not, this file should be manually
	// created.
	//
	// TODO: Pending implementation in GH-122.
	//
	// This will likely require reworking the configuration struct and other
	// config load behavior.
	// defaultTOMLConfigFileName string = "config.toml"
)

const (
	// netTypeTCPAuto is a custom keyword indicating that either of IPv4 or
	// IPv6 is an acceptable network type.
	netTypeTCPAuto string = "auto"

	// netTypeTCP4 indicates that IPv4 network connections are required.
	netTypeTCP4 string = "tcp4"

	// netTypeTCP6 indicates that IPv6 network connections are required
	netTypeTCP6 string = "tcp6"
)

// these keys are found in the `DEFAULT` section of the INI file.
const (
	iniDefaultServerNameKeyName string = "server_name"
	iniDefaultServerPortKeyName string = "server_port"
)

// these keys are found in the other (unique) sections in the INI file.
const (
	iniUsernameKeyName string = "username"
	iniPasswordKeyName string = "password"
	iniFoldersKeyName  string = "folders"
)

// default permissions granting owner full access, deny access to all others
const (
	defaultDirectoryPerms os.FileMode = 0700
	defaultFilePerms      os.FileMode = 0600
)

// logFilenameTemplate is used as a filename template when generating new log
// files. The log file created using this template is appended to by the
// logger throughout application execution. In order to help prevent
// repeatedly writing to the same log file, the current date/time when first
// writing to the file is used as part of the filename itself.
const logFilenameTemplate string = myAppName + "-%s.txt"

// logFilenameDateLayout is the time formatting layout for generated
// reports.
const logFilenameDateLayout string = "2006-01-02-15_04"
