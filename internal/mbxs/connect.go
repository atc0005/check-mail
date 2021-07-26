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

// openConnection receives a list of IP Addresses and returns a client
// connection for the first successful connection attempt. An error is
// returned instead if one occurs.
func openConnection(addrs []string, port int, dialer Dialer, tlsConfig *tls.Config, logger zerolog.Logger) (*client.Client, error) {

	if len(addrs) < 1 {
		errMsg := "empty list of IP Addresses received"

		logger.Error().Msg(errMsg)

		return nil, fmt.Errorf(errMsg)
	}

	var c *client.Client
	var connectErr error

	for _, addr := range addrs {
		logger.Debug().
			Str("ip_address", addr).
			Msg("Connecting to server")

		s := net.JoinHostPort(addr, strconv.Itoa(port))

		// pass in explicitly set TLS config using provided server name, but
		// attempt to connect to specific IP Address returned from earlier
		// lookup. We'll attempt to loop over each available IP Address until
		// we are able to successfully connect to one of them.
		c, connectErr = client.DialWithDialerTLS(&dialer, s, tlsConfig)

		// log override just before checking for an error; this value could be
		// useful in troubleshooting why a connection attempt fails
		if dialer.NetworkTypeUserOverride != "" {
			logger.Debug().
				Str("dialer_network_original_value", dialer.NetworkTypeOriginalValue).
				Str("dialer_network_overridden_value", dialer.NetworkTypeUserOverride).
				Msg("dialer network overridden with user supplied value")
		}

		if connectErr != nil {
			logger.Error().
				Err(connectErr).
				Str("ip_address", addr).
				Msg("error connecting to server")

			continue
		}

		// If no connection errors were received, we can consider the
		// connection attempt a success, clear any previous error and abort
		// attempts to connect to any remaining IP Addresses for the specified
		// server name.
		logger.Info().
			Str("ip_address", addr).
			Msg("Connected to server")

		return c, nil
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
			Msg(errMsg)

		return nil, fmt.Errorf("%s; last error: %w", errMsg, connectErr)
	}

	return c, nil

}

// Connect opens a connection to the specified IMAP server using the specified
// network type, returns a client connection or an error if one occurs.
func Connect(server string, port int, netType string, minTLSVer uint16, logger zerolog.Logger) (*client.Client, error) {

	logger = logger.With().
		Str("hostname", server).
		Str("net_type", netType).
		Logger()

	logger.Debug().Msg("resolving hostname")
	lookupResults, lookupErr := net.LookupHost(server)
	if lookupErr != nil {
		errMsg := "error resolving hostname " + server
		logger.Error().Err(lookupErr).Msg(errMsg)

		return nil, fmt.Errorf(
			"error resolving hostname %s: %w",
			server,
			lookupErr,
		)
	}

	switch {
	case len(lookupResults) < 1:
		errMsg := fmt.Sprintf(
			"failed to resolve hostname %s to IP Addresses",
			server,
		)

		logger.Error().
			Msg(errMsg)

		return nil, fmt.Errorf(errMsg)

	default:
		logger.Debug().
			Int("count", len(lookupResults)).
			Str("ips", strings.Join(lookupResults, ", ")).
			Msg("successfully resolved IP Addresses for hostname")

	}

	addrs := make([]string, 0, len(lookupResults))
	ips := make([]net.IP, 0, len(lookupResults))

	logger.Debug().Msg("converting DNS lookup results to net.IP values for net type validation")
	for i := range lookupResults {
		ip := net.ParseIP(lookupResults[i])
		if ip == nil {
			return nil, fmt.Errorf(
				"error parsing %s as an IP Address",
				lookupResults[i],
			)
		}
		ips = append(ips, ip)
	}

	switch {
	case len(ips) < 1:
		errMsg := fmt.Sprintf(
			"failed to to convert DNS lookup results to net.IP values after receiving %d DNS lookup results ([%s])",
			len(lookupResults),
			strings.Join(lookupResults, ", "),
		)

		logger.Error().Msg(errMsg)

		return nil, fmt.Errorf(errMsg)

	default:
		logger.Debug().Msg("successfully converted DNS lookup results to net.IP values")
	}

	var dialer Dialer

	// Flag validation ensures that we see valid named networks as supported
	// by the `net` stdlib package, along with the "auto" keyword. Here we pay
	// attention to only the valid named networks. Since we're working with
	// user specified keywords, we compare case-insensitively.
	switch strings.ToLower(netType) {
	case NetTypeTCP4:
		logger.Debug().Msg("user opted for IPv4-only connectivity, gathering only IPv4 addresses")
		for i := range ips {
			if ips[i].To4() != nil {
				addrs = append(addrs, ips[i].String())
			}
		}
		dialer.NetworkTypeUserOverride = NetTypeTCP4

	case NetTypeTCP6:
		logger.Debug().Msg("user opted for IPv6-only connectivity, gathering only IPv6 addresses")
		for i := range ips {
			if ips[i].To4() == nil {
				// if earlier attempts to parse the IP Address succeeded, but
				// this is not considered an IPv4 address, we will consider it
				// a valid IPv6 address.
				addrs = append(addrs, ips[i].String())
			}
		}

		dialer.NetworkTypeUserOverride = NetTypeTCP6

	// either of IPv4 or IPv6 is acceptable
	default:
		logger.Debug().Msg("auto behavior enabled, gathering all addresses")
		addrs = lookupResults
	}

	// No IPs remain after filtering against IPv4-only or IPv6-only
	// requirement.
	switch {
	case len(addrs) < 1:
		errMsg := fmt.Sprintf(
			"failed to gather IP Addresses for connection attempts after receiving and parsing %d DNS lookup results ([%s])",
			len(lookupResults),
			strings.Join(lookupResults, ", "),
		)

		logger.Error().Msg(errMsg)

		return nil, fmt.Errorf(errMsg)

	default:
		logger.Debug().
			Int("count", len(addrs)).
			Str("ips", strings.Join(addrs, ", ")).
			Msg("successfully gathered IP Addresses for connection attempts")
	}

	// #nosec G402; allow user to choose minimum TLS version, fallback to a
	// secure default
	tlsConfig := &tls.Config{
		ServerName: server,
		MinVersion: minTLSVer,
	}

	c, connectErr := openConnection(addrs, port, dialer, tlsConfig, logger)
	if connectErr != nil {
		return nil, connectErr
	}

	if c == nil {
		return nil, fmt.Errorf(
			"failed to create client connection to %s using any of IPs %s",
			server,
			strings.Join(addrs, ", "),
		)
	}

	return c, nil

}

// Login uses the provided client connection and credentials to login to the
// remote server.
func Login(client *client.Client, username string, password string, logger zerolog.Logger) error {

	if client == nil {
		errMsg := fmt.Sprintf(
			"invalid (nil) client received while attempting login for account %s",
			username,
		)

		logger.Error().Msg(errMsg)

		return fmt.Errorf(errMsg)
	}

	logger.Debug().Msg("Logging in")
	if err := client.Login(username, password); err != nil {
		errMsg := "login error occurred"
		logger.Error().Err(err).Msg(errMsg)

		return fmt.Errorf("%s: %w", errMsg, err)
	}
	logger.Info().Msg("Logged in")

	return nil

}
