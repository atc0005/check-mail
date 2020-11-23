// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "flag"

func (c *Config) handleFlagsConfig(acceptConfigFile bool) {

	var account MailAccount

	// shared flags
	flag.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)
	flag.StringVar(&c.LoggingLevel, "log-level", defaultLoggingLevel, loggingLevelFlagHelp)

	// currently only applies to Nagios plugin
	if !acceptConfigFile {
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

	// if CLI-provided values were given then record those as an entry in the
	// list
	if !acceptConfigFile {
		c.Accounts = append(c.Accounts, account)
	}

}
