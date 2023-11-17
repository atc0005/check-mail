// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"os"
)

// handleFlagsConfig handles toggling the exposure of specific configuration
// flags to the user. This behavior is controlled via the specified
// application type as set by each cmd. Based on the application's specified
// type, a smaller subset of flags specific to each type are exposed along
// with a set common to all application types.
func (c *Config) handleFlagsConfig(appType AppType) error {

	if c == nil {
		return fmt.Errorf(
			"nil configuration, cannot process flags: %w",
			ErrConfigNotInitialized,
		)
	}

	var account MailAccount

	// shared flags
	c.flagSet.BoolVar(&c.ShowHelp, HelpFlagShort, defaultHelp, helpFlagHelp+shorthandFlagSuffix)
	c.flagSet.BoolVar(&c.ShowHelp, HelpFlagLong, defaultHelp, helpFlagHelp)

	c.flagSet.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)
	c.flagSet.StringVar(&c.LoggingLevel, "log-level", defaultLoggingLevel, loggingLevelFlagHelp)

	if appType.ReporterIMAPMailbox {
		c.flagSet.StringVar(&c.minTLSVersion, "min-tls", defaultMinTLSVersion, minTLSVersionFlagHelp)
		c.flagSet.StringVar(&c.NetworkType, "net-type", defaultNetworkType, networkTypeFlagHelp)
		c.flagSet.StringVar(&c.ConfigFile, "config-file", defaultINIConfigFileName, iniConfigFileFlagHelp)
		c.flagSet.StringVar(&c.ReportFileOutputDir, "report-file-dir", defaultReportFileOutputDir, reportFileOutputDirFlagHelp)
		c.flagSet.StringVar(&c.LogFileOutputDir, "log-file-dir", defaultLogFileOutputDir, logFileOutputDirFlagHelp)
	}

	if appType.InspectorIMAPCaps {
		c.flagSet.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		c.flagSet.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
		c.flagSet.StringVar(&c.minTLSVersion, "min-tls", defaultMinTLSVersion, minTLSVersionFlagHelp)
		c.flagSet.StringVar(&c.NetworkType, "net-type", defaultNetworkType, networkTypeFlagHelp)
	}

	if appType.FetcherOAuth2TokenFromAuthServer {
		c.flagSet.Var(&c.FetcherOAuth2TokenSettings.Scopes, "scopes", scopesFlagHelp)
		c.flagSet.StringVar(&c.FetcherOAuth2TokenSettings.ClientID, "client-id", defaultClientID, clientIDFlagHelp)
		c.flagSet.StringVar(&c.FetcherOAuth2TokenSettings.ClientSecret, "client-secret", defaultClientSecret, clientSecretFlagHelp)
		c.flagSet.StringVar(&c.FetcherOAuth2TokenSettings.TokenURL, "token-url", defaultTokenURL, tokenURLFlagHelp)
		c.flagSet.BoolVar(&c.FetcherOAuth2TokenSettings.EmitTokenAsJSON, "json-output", defaultEmitTokenAsJSON, emitTokenAsJSONFlagHelp)
		c.flagSet.StringVar(&c.FetcherOAuth2TokenSettings.Filename, "filename", defaultTokenFilename, tokenFilenameFlagHelp)
		c.flagSet.IntVar(&c.FetcherOAuth2TokenSettings.RetrievalAttempts, "max-attempts", defaultTokenRetrievalAttempts, tokenRetrievalAttemptsFlagHelp)
	}

	if appType.FetcherOAuth2TokenFromCache {
		c.flagSet.StringVar(&c.FetcherOAuth2TokenSettings.Filename, "filename", defaultTokenFilename, tokenFilenameFlagHelp)
	}

	if appType.PluginIMAPMailboxBasicAuth {

		// Indicate what validation logic should be applied for this set of
		// flags.
		account.AuthType = AuthTypeBasic

		c.flagSet.Var(&account.Folders, "folders", foldersFlagHelp)
		c.flagSet.StringVar(&account.Username, "username", defaultUsername, usernameFlagHelp)
		c.flagSet.StringVar(&account.Password, "password", defaultPassword, passwordFlagHelp)
		c.flagSet.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		c.flagSet.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
		c.flagSet.BoolVar(&c.EmitBranding, "branding", defaultEmitBranding, emitBrandingFlagHelp)
		c.flagSet.StringVar(&c.minTLSVersion, "min-tls", defaultMinTLSVersion, minTLSVersionFlagHelp)
		c.flagSet.StringVar(&c.NetworkType, "net-type", defaultNetworkType, networkTypeFlagHelp)
	}

	if appType.PluginIMAPMailboxOAuth2 {

		// Indicate what validation logic should be applied for this set of
		// flags.
		account.AuthType = AuthTypeOAuth2ClientCreds

		// Common plugin flags
		c.flagSet.Var(&account.Folders, "folders", foldersFlagHelp)
		c.flagSet.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		c.flagSet.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
		c.flagSet.BoolVar(&c.EmitBranding, "branding", defaultEmitBranding, emitBrandingFlagHelp)
		c.flagSet.StringVar(&c.minTLSVersion, "min-tls", defaultMinTLSVersion, minTLSVersionFlagHelp)
		c.flagSet.StringVar(&c.NetworkType, "net-type", defaultNetworkType, networkTypeFlagHelp)

		// OAuth2 flags
		c.flagSet.Var(&account.OAuth2Settings.Scopes, "scopes", scopesFlagHelp)
		c.flagSet.StringVar(&account.OAuth2Settings.ClientID, "client-id", defaultClientID, clientIDFlagHelp)
		c.flagSet.StringVar(&account.OAuth2Settings.ClientSecret, "client-secret", defaultClientSecret, clientSecretFlagHelp)
		c.flagSet.StringVar(&account.OAuth2Settings.SharedMailbox, "shared-mailbox", defaultSharedMailbox, sharedMailboxFlagHelp)
		c.flagSet.StringVar(&account.OAuth2Settings.TokenURL, "token-url", defaultTokenURL, tokenURLFlagHelp)

	}

	// Allow our function to override the default Help output.
	//
	// Override default of stderr as destination for help output. This allows
	// Nagios XI and similar monitoring systems to call plugins with the
	// `--help` flag and have it display within the Admin web UI.
	c.flagSet.Usage = Usage(c.flagSet, os.Stdout)

	// parse flag definitions from the argument list
	if err := c.flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	// For all app types other than the Reporter app we need to save any
	// configured account details provided via CLI; the Reporter app receives
	// all account details via configuration file.
	if !appType.ReporterIMAPMailbox {
		c.Accounts = append(c.Accounts, account)
	}

	return nil

}
