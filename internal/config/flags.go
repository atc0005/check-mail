// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "flag"

func (c *Config) handleFlagsConfig() {

	flag.Var(&c.Folders, "folders", foldersFlagHelp)
	flag.StringVar(&c.Username, "username", defaultUsername, usernameFlagHelp)
	flag.StringVar(&c.Password, "password", defaultPassword, passwordFlagHelp)
	flag.StringVar(&c.Server, "server", defaultServer, serverFlagHelp)
	flag.IntVar(&c.Port, "port", defaultPort, portFlagHelp)
	flag.StringVar(&c.LoggingLevel, "log-level", defaultLoggingLevel, loggingLevelFlagHelp)
	flag.BoolVar(&c.EmitBranding, "branding", defaultEmitBranding, emitBrandingFlagHelp)
	flag.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

}
