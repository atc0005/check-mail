// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

//go:generate go-winres make --product-version=git-tag --file-version=git-tag

package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/check-mail/internal/mbxs"
	"github.com/rs/zerolog"
)

func main() {

	// Setup configuration by parsing user-provided flags
	cfg, cfgErr := config.New(config.AppType{InspectorIMAPCaps: true})
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())

		return

	case errors.Is(cfgErr, config.ErrHelpRequested):
		fmt.Println(cfg.Help())

		return

	case cfgErr != nil:

		// We make some assumptions when setting up our logger as we do not
		// have a working configuration based on sysadmin-specified choices.
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		logger := zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()

		logger.Err(cfgErr).Msg("Error initializing application")

		return
	}

	logger := cfg.Log.With().Logger()

	// We're reusing the common "Accounts" field (and flags) in order to
	// obtain specified server and port values vs adding standalone fields and
	// flags; this is effectively a loop of 1 iteration (at least for now). At
	// some point we may expand the scope of this tool to handle evaluating
	// mail servers for a collection of accounts, so using a workflow intended
	// for collections is probably appropriate.
	for _, account := range cfg.Accounts {

		// Open connection to IMAP server
		c, err := mbxs.Connect(account.Server, account.Port, cfg.NetworkType, cfg.MinTLSVersion(), logger)
		if err != nil {
			logger.Error().Err(err).Msg("error connecting to server")
			os.Exit(1)
		}
		logger.Info().Msg("Connection established to server")

		// Enable client network command/response logging if global logging
		// level indicates user wishes to see verbose details.
		if zerolog.GlobalLevel() == zerolog.DebugLevel ||
			zerolog.GlobalLevel() == zerolog.TraceLevel {
			c.SetDebug(&logger)
		}

		logger.Info().Msg("Gathering pre-login capabilities")
		capabilities, err := c.Capability()
		if err != nil {
			logger.Error().Err(err).Msg("Unable to list server capabilities")
			os.Exit(1)
		}

		caps := make([]string, 0, len(capabilities))
		for k, v := range capabilities {
			if v {
				caps = append(caps, k)
			}
		}

		sort.Strings(caps)
		// logger.Info().Msgf("Capabilities: %v", caps)
		for _, capability := range caps {
			logger.Info().Msgf("Capability: %v", capability)
		}

		logger.Debug().Msg("Closing connection to server")
		if err := c.Logout(); err != nil {
			logger.Error().Err(err).Msg("failed to close connection to server")
			os.Exit(1)
		}
		logger.Info().Msg("Connection to server closed")

	}
}
