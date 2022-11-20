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

	// currently only applies to list-emails app, don't expose to Nagios plugin
	if appType.ReporterIMAPMailboxBasicAuth {
		flag.StringVar(&c.ConfigFile, "config-file", defaultINIConfigFileName, iniConfigFileFlagHelp)
		flag.StringVar(&c.ReportFileOutputDir, "report-file-dir", defaultReportFileOutputDir, reportFileOutputDirFlagHelp)
		flag.StringVar(&c.LogFileOutputDir, "log-file-dir", defaultLogFileOutputDir, logFileOutputDirFlagHelp)
	}

	// currently only applies to Nagios plugin
	if appType.PluginIMAPMailboxBasicAuth {
		flag.Var(&account.Folders, "folders", foldersFlagHelp)
		flag.StringVar(&account.Username, "username", defaultUsername, usernameFlagHelp)
		flag.StringVar(&account.Password, "password", defaultPassword, passwordFlagHelp)
		flag.StringVar(&account.Server, "server", defaultServer, serverFlagHelp)
		flag.IntVar(&account.Port, "port", defaultPort, portFlagHelp)
		flag.BoolVar(&c.EmitBranding, "branding", defaultEmitBranding, emitBrandingFlagHelp)
	}

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

	// For all app types other than the Reporter app we need to save any
	// configured account details provided via CLI; the Reporter app receives
	// all account details via configuration file.
	if !appType.ReporterIMAPMailboxBasicAuth {
		c.Accounts = append(c.Accounts, account)
	}

}
