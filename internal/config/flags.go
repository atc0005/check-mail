// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "flag"

// handleFlagsConfig handles toggling the exposure of specific configuration
// flags to the user. This behavior is controlled via the specified
// application type as set by each cmd. Based on the application's specified
// type, a smaller subset of flags specific to each type are exposed along
// with a set common to all application types.
func (c *Config) handleFlagsConfig(appType AppType) {

	var account MailAccount

	// shared flags
	flag.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)
	flag.StringVar(&c.LoggingLevel, "log-level", defaultLoggingLevel, loggingLevelFlagHelp)
	flag.StringVar(&c.NetworkType, "net-type", defaultNetworkType, networkTypeFlagHelp)
	flag.StringVar(&c.minTLSVersion, "min-tls", defaultMinTLSVersion, minTLSVersionFlagHelp)

	// Only applies to Reporter app
	if appType.ReporterIMAPMailbox {
		flag.StringVar(&c.ConfigFile, "config-file", defaultINIConfigFileName, iniConfigFileFlagHelp)
		flag.StringVar(&c.ReportFileOutputDir, "report-file-dir", defaultReportFileOutputDir, reportFileOutputDirFlagHelp)
		flag.StringVar(&c.LogFileOutputDir, "log-file-dir", defaultLogFileOutputDir, logFileOutputDirFlagHelp)
	}

	// Inspector app
	if appType.InspectorIMAPCaps {
		flag.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		flag.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
	}

	// Basic Auth Plugin
	if appType.PluginIMAPMailboxBasicAuth {

		// Indicate what validation logic should be applied for this set of
		// flags.
		account.AuthType = AuthTypeBasic

		flag.Var(&account.Folders, "folders", foldersFlagHelp)
		flag.StringVar(&account.Username, "username", defaultUsername, usernameFlagHelp)
		flag.StringVar(&account.Password, "password", defaultPassword, passwordFlagHelp)
		flag.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		flag.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
		flag.BoolVar(&c.EmitBranding, "branding", defaultEmitBranding, emitBrandingFlagHelp)
	}

	// OAuth2 Client Credentials flow plugin
	if appType.PluginIMAPMailboxOAuth2 {

		// Indicate what validation logic should be applied for this set of
		// flags.
		account.AuthType = AuthTypeOAuth2ClientCreds

		// Common plugin flags
		flag.Var(&account.Folders, "folders", foldersFlagHelp)
		flag.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		flag.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
		flag.BoolVar(&c.EmitBranding, "branding", defaultEmitBranding, emitBrandingFlagHelp)

		// OAuth2 flags
		flag.Var(&account.OAuth2Settings.Scopes, "scopes", scopesFlagHelp)
		flag.StringVar(&account.OAuth2Settings.ClientID, "client-id", defaultClientID, clientIDFlagHelp)
		flag.StringVar(&account.OAuth2Settings.ClientSecret, "client-secret", defaultClientSecret, clientSecretFlagHelp)
		flag.StringVar(&account.OAuth2Settings.SharedMailbox, "shared-mailbox", defaultSharedMailbox, sharedMailboxFlagHelp)
		flag.StringVar(&account.OAuth2Settings.TokenURL, "token-url", defaultTokenURL, tokenURLFlagHelp)

	}

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

	// For all app types other than the Reporter app we need to save any
	// configured account details provided via CLI; the Reporter app receives
	// all account details via configuration file.
	if !appType.ReporterIMAPMailbox {
		c.Accounts = append(c.Accounts, account)
	}

}
