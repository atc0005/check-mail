// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"os"
	"time"
)

// Flag names. Exported so that they can be referenced from tests.
const (
	HelpFlagLong  string = "help"
	HelpFlagShort string = "h"
)

// Shared flag help text
const (
	foldersFlagHelp       string = "Folders or IMAP \"mailboxes\" to check for mail. This value is provided as a comma-separated list."
	serverFlagHelp        string = "The fully-qualified domain name of the remote mail server."
	portFlagHelp          string = "TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections."
	networkTypeFlagHelp   string = "Limits network connections to remote mail servers to one of tcp4 (IPv4-only), tcp6 (IPv6-only) or auto (either)."
	minTLSVersionFlagHelp string = "Limits version of TLS used for connections to remote mail servers to one of tls10 (TLS v1.0), tls11, tls12 or tls13 (TLS v1.3)."
	loggingLevelFlagHelp  string = "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace."
	emitBrandingFlagHelp  string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	helpFlagHelp          string = "Emit this help text"
	versionFlagHelp       string = "Whether to display application version and then immediately exit application."
)

// PluginIMAPMailboxBasicAuth flag help text
const (
	usernameFlagHelp string = "The account used to login to the remote mail server using Basic Auth. This is often in the form of an email address."
	passwordFlagHelp string = "The remote mail server account password. Used for Basic Auth."
)

// PluginIMAPMailboxOauth2 flag help text
const (
	clientIDFlagHelp      string = "Application (client) ID created during app registration."
	clientSecretFlagHelp  string = "Client secret (aka, \"app\" password)."
	scopesFlagHelp        string = "One or more scopes requested from the authorization server. E.g., \"https://outlook.office365.com/.default\" for O365."
	sharedMailboxFlagHelp string = "Email account that is to be accessed using client ID & secret values. Usually a shared mailbox among a team."

	// False-positive gosec linter warning
	//nolint
	tokenURLFlagHelp string = "The OAuth2 provider's token endpoint URL. E.g., \"https://accounts.google.com/o/oauth2/token\" for Google. See example INI file for O365 example."
)

// Reporter flag help text
const (
	iniConfigFileFlagHelp       string = "Full path to the INI-formatted configuration file used by this application. See the accounts.example.ini files under contrib/list-emails directory for a starter template. Copy to accounts.ini, update with applicable information and place in a directory of your choice. If this file is found in your current working directory you need not use this flag."
	reportFileOutputDirFlagHelp string = "Full path to the directory where email summary report files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist."
	logFileOutputDirFlagHelp    string = "Full path to the directory where log files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist."
)

// Fetcher flag help text
const (
	emitTokenAsJSONFlagHelp string = "Emit retrieved token in JSON format. Defaults to emitting the access token field from retrieved payload."
	tokenFilenameFlagHelp   string = "Save retrieved token to specified file. Emitted to standard out (stdout) if not specified."

	// False-positive gosec linter warning
	//nolint
	tokenRetrievalAttemptsFlagHelp string = "Max token retrieval attempts."
)

// shorthandFlagSuffix is appended to short flag help text to emphasize that
// the flag is a shorthand version of a longer flag.
const shorthandFlagSuffix = " (shorthand)"

// Default flag settings if not overridden by user input
const (
	defaultHelp                  bool   = false
	defaultLoggingLevel          string = "info"
	defaultEmitBranding          bool   = false
	defaultPort                  int    = 993
	defaultServer                string = ""
	defaultPassword              string = ""
	defaultUsername              string = ""
	defaultClientID              string = ""
	defaultClientSecret          string = ""
	defaultSharedMailbox         string = ""
	defaultTokenURL              string = ""
	defaultNetworkType           string = netTypeTCPAuto
	defaultMinTLSVersion         string = minTLSVersion12
	defaultDisplayVersionAndExit bool   = false
	defaultEmitTokenAsJSON       bool   = false
	defaultTokenFilename         string = ""

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

	defaultAccountProcessDelay time.Duration = time.Second * 5

	defaultTokenRetrievalAttempts int = 3
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

// TLS keywords used to map to TLS versions in the tls stdlib package.
// https://golang.org/pkg/crypto/tls/#pkg-constants
const (
	minTLSVersion10 string = "tls10"
	minTLSVersion11 string = "tls11"
	minTLSVersion12 string = "tls12"
	minTLSVersion13 string = "tls13"
)

// These keys are found in the `DEFAULT` section of the INI file. The design
// is that these values are server-specific while the other keys are
// account/mailbox specific.
//
// The INI config file layout is one server (potentially with one
// client/application ID) and one or more associated accounts/mailboxes.
const (
	iniDefaultAuthTypeKeyName         string = "auth_type"
	iniDefaultServerNameKeyName       string = "server_name"
	iniDefaultServerPortKeyName       string = "server_port"
	iniDefaultClientIDKeyName         string = "client_id"
	iniDefaultClientSecretKeyName     string = "client_secret"
	iniDefaultScopesKeyName           string = "scopes"
	iniDefaultEndpointTokenURLKeyName string = "endpoint_token_url"
)

// These keys are found in the other (unique) sections in the INI file. If
// account is specified then username will not be.
const (
	iniUsernameKeyName      string = "username"
	iniSharedMailboxKeyName string = "shared_mailbox"
	iniPasswordKeyName      string = "password"
	iniFoldersKeyName       string = "folders"
)

// Supported authentication types used by applications in this project.
//
// These are also the only valid values for the auth_type INI config file key.
const (

	// AuthTypeBasic indicates Basic Authentication (username/password).
	AuthTypeBasic string = "basic"

	// AuthTypeOAuth2ClientCreds indicates OAuth2 Client Credentials flow.
	AuthTypeOAuth2ClientCreds string = "oauth2"
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
