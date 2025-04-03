// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

//go:generate go-winres make --product-version=git-tag --file-version=git-tag

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"

	"github.com/atc0005/check-mail/internal/config"
)

type exitStatus struct {
	Code int
	Err  error
}

func main() {

	ctx := context.Background()

	var appExitStatus exitStatus

	defer func(appExitStatus *exitStatus) {
		if appExitStatus.Err != nil {
			fmt.Println(appExitStatus.Err.Error())
		}
		os.Exit(appExitStatus.Code)
	}(&appExitStatus)

	// Setup configuration by parsing user-provided flags
	cfg, cfgErr := config.New(config.AppType{ReporterIMAPMailbox: true})
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())
		appExitStatus.Code = 0

		return

	case errors.Is(cfgErr, config.ErrHelpRequested):
		fmt.Println(cfg.Help())
		appExitStatus.Code = 0

		return

	case cfgErr != nil:
		// We're using the standalone Err function from rs/zerolog/log as we
		// do not have a working configuration with a preconfigured logger
		// attached. By default, this message will go to stderr and should
		// be decipherable by the user.
		zlog.Err(cfgErr).Msg("Error initializing application")

		appExitStatus.Code = 1

		return
	}

	// #nosec G307
	// Believed to be a false-positive from recent gosec release
	// https://github.com/securego/gosec/issues/714
	defer func(filename string) {
		if err := cfg.LogFileHandle.Close(); err != nil {
			// Ignore "file already closed" errors
			if !errors.Is(err, os.ErrClosed) {
				// We're using the standalone Err function from rs/zerolog/log
				// as we have the main logger set to write to a log file,
				// which we just failed to close. By default, this message
				// will go to stderr and should be decipherable by the user.
				zlog.Error().
					Err(err).
					Str("filename", filename).
					Msg("failed to close file")
			}
		}
	}(cfg.LogFileHandle.Name())

	// loop over accounts
	for i, account := range cfg.Accounts {
		// Building with `go build -gcflags=all=-d=loopvar=2` identified this
		// loop as compiling differently with Go 1.22 (per-iteration) loop
		// semantics.
		//
		// As a workaround, we create a new variable for each iteration to
		// work around potential issues with Go versions prior to Go 1.22.
		//
		// NOTE: Not needed as of Go 1.22.
		//
		// account := account

		fmt.Println("Checking account:", account.Name)

		logger := cfg.Log.With().
			Str("auth_type", account.AuthType).
			Str("server", account.Server).
			Int("port", account.Port).
			Str("folders_to_check", account.Folders.String()).
			Logger()

		if err := processAccount(ctx, account, cfg, logger); err != nil {
			appExitStatus.Err = err
			appExitStatus.Code = 1
			return
		}

		if i+1 < len(cfg.Accounts) {
			// Delay processing the next account (unless we've processed them
			// all) in an attempt to prevent encountering the "User is
			// authenticated but not connected" error that is believed to
			// occur when remote connections limit is exceeded.
			time.Sleep(cfg.AccountProcessDelay())
		}

	}

	fmt.Printf(
		"OK: Successfully generated reports for accounts: %s\n",
		strings.Join(cfg.AccountNames(), ", "),
	)

	// implied return here :)

}
