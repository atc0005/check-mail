// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog"
)

// Connect opens a connection to the specified IMAP server, returns a client
// connection.
func Connect(server string, port int, logger zerolog.Logger) (*client.Client, error) {

	logger.Debug().Msg("resolving hostname")
	addrs, lookupErr := net.LookupHost(server)
	if lookupErr != nil {
		errMsg := "error resolving hostname " + server
		logger.Error().Err(lookupErr).Msg(errMsg)

		return nil, fmt.Errorf(
			"error resolving hostname %s: %w",
			server,
			lookupErr,
		)
	}

	var c *client.Client
	var connectErr error
	tlsConfig := &tls.Config{
		ServerName: server,
		// NOTE: Explicitly setting minimum TLS version to resolve `G402: TLS
		// MinVersion too low. (gosec)`
		//
		// TODO: Revisit as part of GH-169
		MinVersion: tls.VersionTLS12,
	}
	dialer := &net.Dialer{}

	for _, addr := range addrs {
		logger.Debug().
			Str("ip_address", addr).
			Str("hostname", server).
			Msg("Connecting to server")

		s := net.JoinHostPort(addr, strconv.Itoa(port))

		// pass in explicitly set TLS config using provided server name, but
		// attempt to connect to specific IP Address returned from earlier
		// lookup. We'll attempt to loop over each available IP Address until
		// we are able to successfully connect to one of them.
		c, connectErr = client.DialWithDialerTLS(dialer, s, tlsConfig)
		if connectErr != nil {
			logger.Error().
				Err(connectErr).
				Str("ip_address", addr).
				Str("hostname", server).
				Msg("error connecting to server")

			continue
		}

		// If no connection errors were received, we can consider the
		// connection attempt a success, clear any previous error and abort
		// attempts to connect to any remaining IP Addresses for the specified
		// server name.
		logger.Info().
			Str("ip_address", addr).
			Str("hostname", server).
			Msg("Connected to server")
		connectErr = nil
		break
	}

	// If all connection attempts failed, report the last connection error.
	// Log all failed IP Addresses for review.
	if connectErr != nil {
		errMsg := fmt.Sprintf(
			"failed to connect to server using any of %d IP Addresses (%s)",
			len(addrs),
			strings.Join(addrs, ", "),
		)
		logger.Error().
			Err(connectErr).
			Str("failed_ip_addresses", strings.Join(addrs, ", ")).
			Str("hostname", server).
			Msg(errMsg)

		return nil, fmt.Errorf("%s: %w", errMsg, connectErr)
	}

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
