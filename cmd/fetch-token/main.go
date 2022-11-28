// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

//go:generate go-winres make --product-version=git-tag --file-version=git-tag

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/check-mail/internal/oauth2"
	"github.com/rs/zerolog"
)

func main() {

	ctx := context.Background()

	// Setup configuration by parsing user-provided flags
	cfg, cfgErr := config.New(config.AppType{FetcherOAuth2TokenFromAuthServer: true})
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

	var logger zerolog.Logger
	switch {
	case cfg.FetcherOAuth2TokenSettings.Filename != "":
		logger = cfg.Log.With().
			Str("filename", cfg.FetcherOAuth2TokenSettings.Filename).
			Logger()
	default:
		logger = cfg.Log.With().Logger()
	}

	logger.Debug().Msg("Application configuration initialized")

	logger.Debug().Msg("Fetching Client Credentials token")
	token, err := oauth2.GetClientCredentialsToken(
		ctx,
		cfg.FetcherOAuth2TokenSettings.ClientID,
		cfg.FetcherOAuth2TokenSettings.ClientSecret,
		cfg.FetcherOAuth2TokenSettings.Scopes,
		cfg.FetcherOAuth2TokenSettings.TokenURL,
		cfg.RetrievalAttempts(),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve token")
		os.Exit(1)
	}
	logger.Debug().
		Str("token_expiration", token.Expiry.Format(time.RFC3339)).
		Str("token_type", token.Type()).
		Msg("Token retrieved")

	var data []byte
	var emittedAsJSON bool
	switch {

	case cfg.FetcherOAuth2TokenSettings.EmitTokenAsJSON:
		var err error
		data, err = json.MarshalIndent(token, "", "\t")
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to marshal token to JSON format")
			os.Exit(1)
		}
		logger.Debug().Msg("Successfully converted token to JSON")

		emittedAsJSON = true

	default:
		logger.Debug().Msg("Retaining access token as plaintext value")
		data = []byte(token.AccessToken)
		emittedAsJSON = false
	}

	switch {
	case cfg.FetcherOAuth2TokenSettings.Filename != "":
		err := os.WriteFile(filepath.Clean(cfg.FetcherOAuth2TokenSettings.Filename), data, 0600)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to write data to output file")
			os.Exit(1)
		}

		logger.Debug().Msg("Successfully wrote data to file")

	default:
		n, err := os.Stdout.Write(data)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to write data to stdout")
			os.Exit(1)
		}
		logger.Debug().
			Int("bytes_written", n).
			Bool("emitted_as_json", emittedAsJSON).
			Msg("Emitted retrieved token")
	}

}
