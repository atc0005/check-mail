// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import (
	"fmt"

	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog"
)

// Connect opens a connection to the specified IMAP server, returns a client
// connection.
func Connect(server string, port int, logger zerolog.Logger) (*client.Client, error) {

	s := fmt.Sprintf("%s:%d", server, port)

	logger.Debug().Msg("connecting to remote server")
	c, err := client.DialTLS(s, nil)
	if err != nil {
		errMsg := "error connecting to server"
		logger.Error().Err(err).Msgf(errMsg)

		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	logger.Info().Msg("Connected")

	return c, nil

}

// Login uses the provided client connection and credentials to login to the
// remote server.
func Login(client *client.Client, username string, password string, logger zerolog.Logger) error {

	logger.Debug().Msg("Logging in")
	if err := client.Login(username, password); err != nil {
		errMsg := "login error occurred"
		logger.Error().Err(err).Msg(errMsg)

		return fmt.Errorf("%s: %w", errMsg, err)
	}
	logger.Info().Msg("Logged in")

	return nil

}
