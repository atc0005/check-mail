// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Updated via Makefile builds. Setting placeholder value here so that
// something resembling a version string will be provided for non-Makefile
// builds.
var version = "x.y.z"

const myAppName string = "check-mail"
const myAppURL string = "https://github.com/atc0005/check-mail"

var (
	// ErrVersionRequested indicates that the user requested application
	// version information
	ErrVersionRequested = errors.New("version information requested")

	// ErrAppTypeNotSpecified indicates that a tool in this project failed to
	// specify a valid application type.
	ErrAppTypeNotSpecified = errors.New("valid app type not specified")

	// ErrInvalidAuthType indicates that an invalid or unsupported
	// authentication type was specified.
	ErrInvalidAuthType = errors.New("invalid auth type")
)

// AppType represents the type of application that is being
// configured/initialized. Not all application types will use the same
// features and as a result will not accept the same flags. Unless noted
// otherwise, each of the application types are incompatible with each other,
// though some flags are common to all types.
type AppType struct {

	// PluginIMAPMailboxBasicAuth represents an application used as a
	// monitoring plugin for evaluating IMAP mailboxes.
	//
	// Basic Authentication is used to login.
	PluginIMAPMailboxBasicAuth bool

	// PluginIMAPMailboxOAuth2 represents an application used as a monitoring
	// plugin for evaluating IMAP mailboxes.
	//
	// An OAuth2 flow is used to login.
	PluginIMAPMailboxOAuth2 bool

	// ReporterIMAPMailbox represents an application used for generating
	// reports for specified IMAP mailboxes.
	//
	// Unlike an Inspector application which is focused on testing or
	// gathering specific details for troubleshooting purposes or a monitoring
	// plugin which is intended for providing a severity-based outcome, a
	// Reporter application is intended for gathering information as an
	// overview.
	ReporterIMAPMailbox bool

	// InspectorIMAPCaps represents an application used for one-off or
	// isolated checks of an IMAP server's advertised capabilities.
	//
	// Unlike a monitoring plugin which is focused on specific attributes
	// resulting in a severity-based outcome, an Inspector application is
	// intended for examining a small set of targets for
	// informational/troubleshooting purposes.
	InspectorIMAPCaps bool
}

// OAuth2MailAccountSettings is a collection of OAuth2 settings for a mail
// account that applications in this project interact with.
type OAuth2MailAccountSettings struct {
	// ClientID is the client ID used by the application that asks for
	// authorization. It must be unique across all clients that the
	// authorization server handles. This ID represents the registration
	// information provided by the client.
	//
	// The client identifier is not a secret; it is exposed to the resource
	// owner and MUST NOT be used alone for client authentication. The client
	// identifier is unique to the authorization server.
	// https://datatracker.ietf.org/doc/html/rfc6749#section-2.2
	ClientID string `json:"client_id"`

	// ClientSecret is a secret known only to the application and the
	// authorization server. It can be considered the application's own
	// password. This value is provided upon application authorization.
	ClientSecret string `json:"client_secret"`

	// Scopes is the collection of permissions or "scopes" requested by an
	// application from the authorization server.
	//
	// Scopes let you specify exactly what type of access the application
	// needs. Scopes limit access for OAuth tokens. They do not grant any
	// additional permission beyond that which the user already has.
	//
	// https://www.oauth.com/oauth2-servers/scope/
	// https://docs.github.com/en/developers/apps/building-oauth-apps/scopes-for-oauth-apps
	Scopes multiValueFlag `json:"scope,omitempty"`

	// TenantID is the tenant or customer identifier associated with the
	// OAuth2-enabled service. For example, with Microsoft Office 365 (O365)
	// this value is used to represent the organization subscribed to O365
	// services.
	// TenantID string `json:"-"`

	// SharedMailbox is the email account that is to be accessed by the
	// application using the given client ID, client secret values. This is
	// usually a shared mailbox among a team.
	SharedMailbox string `json:"-"`

	// Token is a valid XOAUTH2 encoded token to use in place of requesting a
	// new token from the authorization server. If specified, Token obviates
	// the need for most other values: ClientID, ClientSecret, TenantID,
	// Mailbox, Username and Password.
	//
	// NOTE: This value is supported by the Plugin application type only.
	//
	// TODO: Does this provide sufficient value? Tradeoff of token reuse (and
	// everything required to save/load it) vs fetching a new token ...
	//
	// Token string `json:"-"`

	// TokenURL is the authority endpoint for token retrieval.
	TokenURL string
}

// MailAccount represents an email account. The values are provided via
// command-line flags or are specified within a configuration file.
type MailAccount struct {
	// Server is the FQDN associated with the IMAP server.
	//
	// Use of IP Addresses is discouraged; TLS is needed for secure
	// communication with IMAP servers and IP Addresses are rarely listed in
	// Subject Alternate Names (SANs) lists for certificates.
	Server string

	// Port is the TCP port associated with the IMAP server. This is usually
	// 993 but may differ for some installations.
	Port int

	// AuthType is either "basic" for Basic Authentication or "oauth2" for the
	// Client Credentials OAuth2 flow. This value acts as a logic switch for
	// the Reporter application.
	AuthType string

	// Folders is a collection of paths associated with an account. This
	// includes paths such as "Inbox", "Junk EMail" or "Trash".
	Folders multiValueFlag

	// Username is usually the full email address associated with an account.
	Username string

	// Password is the plaintext password for the email account.
	Password string

	// OAuth2Settings is a collection of settings specific to OAuth2
	// authentication with the service hosting the email account.
	OAuth2Settings OAuth2MailAccountSettings

	// Name is often the bare username for the email account, but may not be.
	// This is used as the section header within the configuration file.
	//
	// NOTE: As of the v0.4.x releases this field is used exclusively by the
	// reporter tool.
	Name string
}

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

	// EmitBranding controls whether "generated by" text is included at the
	// bottom of application output. This output is included in the Nagios
	// dashboard and notifications. This output may not mix well with branding
	// output from other tools such as atc0005/send2teams which also insert
	// their own branding output.
	EmitBranding bool

	// ShowVersion is a flag indicating whether the user opted to display only
	// the version string and then immediately exit the application.
	ShowVersion bool

	// ConfigFileLoaded is an internal flag indicating whether a user-provided
	// config file was specified *and* loaded, or a config file was
	// automatically detected *and* loaded.
	ConfigFileLoaded bool

	// ConfigFile is the path to the user-provided config file. This config
	// file is not currently used by the check_imap_mailbox plugin provided by
	// this project.
	ConfigFile string

	// ConfigFileUsed is an internal field indicating *what* config file was
	// loaded, be it explicitly specified by the user or automatically
	// detected from a known location.
	ConfigFileUsed string

	// NetworkType indicates whether an attempt should be made to connect to
	// only IPv4, only IPv6 or IMAP servers listening on either of IPv4 or
	// IPv6 addresses ("auto").
	NetworkType string

	// minTLSVersion is the keyword representing the minimum version of TLS
	// supported for encrypted IMAP server connections.
	minTLSVersion string

	// ReportFileOutputDir is the full path to the directory where email
	// summary report files will be generated. Not currently used by the
	// Nagios plugin.
	ReportFileOutputDir string

	// LogFileOutputDir is the full path to the directory where log files will
	// be generated. Not currently used by the Nagios plugin.
	LogFileOutputDir string

	// LogFileHandle is reference to a log file for deferred closure.
	LogFileHandle *os.File

	// LoggingLevel is the supported logging level for this application.
	LoggingLevel string

	// Accounts is the collection of IMAP mail accounts checked by
	// applications provided by this project.
	Accounts []MailAccount

	// Log is an embedded zerolog Logger initialized via config.New().
	Log zerolog.Logger
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

// Usage is a custom override for the default Help text provided by the flag
// package. Here we prepend some additional metadata to the existing output.
func Usage() {

	// Override default of stderr as destination for help output. This allows
	// Nagios XI and similar monitoring systems to call plugins with the
	// `--help` flag and have it display within the Admin web UI.
	flag.CommandLine.SetOutput(os.Stdout)

	fmt.Fprintln(flag.CommandLine.Output(), "\n"+Version()+"\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

// New is a factory function that produces a new Config object based on user
// provided flag and config file values. It is responsible for validating
// user-provided values and initializing the logging settings used by this
// application.
func New(appType AppType) (*Config, error) {
	var config Config

	config.handleFlagsConfig(appType)

	if config.ShowVersion {
		return nil, ErrVersionRequested
	}

	// Initial validation pass using flag values only.
	if err := config.validate(appType); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// initialize logging "early", just as soon as validation is complete so
	// that we can rely on it to debug further configuration init work
	if err := config.setupLogging(appType); err != nil {
		return nil, fmt.Errorf(
			"failed to set logging configuration: %w",
			err,
		)
	}

	if appType.ReporterIMAPMailbox {
		if err := config.load(); err != nil {

			// We log this message in an effort to populate the log file with
			// something useful; an empty log file isn't that helpful if
			// someone needs to debug later what happened (and the person
			// running the application didn't catch the error output).
			errMsg := "failed to load configuration file"
			config.Log.Error().Err(err).Msgf(errMsg)

			return nil, fmt.Errorf("%s: %w", errMsg, err)
		}

		// Final validation pass using flag AND config file values.
		if err := config.validate(appType); err != nil {
			return nil, fmt.Errorf(
				"configuration validation after loading config file failed: %w",
				err,
			)
		}
	}

	// fmt.Printf("config: %+v\n", config)

	return &config, nil

}

// load is a helper function to handle the bulk of the configuration loading
// work for the New constructor function.
func (c *Config) load() error {

	configFiles := make([]string, 0, 3)

	switch {

	// If specified, load user-specified config file.
	case c.ConfigFile != "":

		c.Log.Debug().
			Str("config_file_candidate", c.ConfigFile).
			Msg("Trying to load user-requested config file")

		loadErr := c.loadConfigFile(c.ConfigFile)
		if loadErr != nil {
			return fmt.Errorf(
				"failed to load configuration from file %q: %w",
				c.ConfigFile,
				loadErr,
			)
		}

	// If not explicitly specified, attempt to automatically load a
	// configuration file from known locations preferring to load a local
	// configuration file from the current working directory first.
	case c.ConfigFile == "":

		localINIConfig, localFileErr := c.localConfigFile(defaultINIConfigFileName)
		if localFileErr != nil {
			return fmt.Errorf(
				"failed to construct path to local config file :%w",
				localFileErr,
			)
		}
		configFiles = append(configFiles, localINIConfig)

		userConfigFile, userConfigFileErr := c.userConfigFile(
			myAppName, defaultINIConfigFileName,
		)
		if userConfigFileErr != nil {
			return fmt.Errorf(
				"failed to construct path to user config file :%w",
				userConfigFileErr,
			)
		}
		configFiles = append(configFiles, userConfigFile)

		loadErr := c.loadConfigFile(configFiles...)
		if loadErr != nil {
			return fmt.Errorf(
				"failed to load candidate config files %v: %w",
				configFiles,
				loadErr,
			)
		}
	}

	return nil
}
