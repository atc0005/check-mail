// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

//go:generate go-winres make --product-version=git-tag --file-version=git-tag

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

func main() {

	// Setup configuration by parsing user-provided flags
	cfg, cfgErr := config.New(config.AppType{FetcherOAuth2TokenFromCache: true})
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
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}
		logger := zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()

		logger.Err(cfgErr).Msg("Error initializing application")

		return
	}

	logger := cfg.Log.With().
		Str("filename", cfg.FetcherOAuth2TokenSettings.Filename).
		Logger()

	logger.Debug().Msg("Application configuration initialized")

	logger.Debug().Msg("Fetching Client Credentials token from file")
	data, err := os.ReadFile(cfg.FetcherOAuth2TokenSettings.Filename)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read file contents")
		os.Exit(1)
	}
	logger.Debug().Msg("Successfully read file contents")

	var output []byte
	switch {
	case bytes.Contains(data, []byte("{")):
		logger.Error().Err(err).Msg("File contents appear to be JSON, will attempt to parse as JSON")

		var token oauth2.Token
		if err := json.Unmarshal(data, &token); err != nil {
			logger.Error().Err(err).Msg("Failed to parse file contents as JSON")
			os.Exit(1)
		}
		logger.Debug().Msg("Successfully parsed file contents as JSON")

		if !token.Valid() {
			logger.Error().
				Str("token_expiration", token.Expiry.Format(time.RFC3339)).
				Str("token_type", token.Type()).
				Msg("Token is NOT valid; a new token should be retrieved and cached in file")
			os.Exit(1)
		}

		logger.Debug().
			Str("token_expiration", token.Expiry.Format(time.RFC3339)).
			Str("token_type", token.Type()).
			Msg("Token is valid, retrieving access token value")

		output = []byte(token.AccessToken)

	default:
		logger.Debug().Msg("File contents do not appear to be JSON")
		logger.Debug().Msg("Attempting to parse file contents as plaintext access token")
		output = data
	}

	n, err := os.Stdout.Write(output)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to emit token")
		os.Exit(1)
	}
	logger.Debug().
		Int("bytes_written", n).
		Msg("Emitted retrieved token")

}
