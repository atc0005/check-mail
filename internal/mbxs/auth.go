// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atc0005/check-mail/internal/sasl"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	// ErrRequiredAuthMechanismUnsupported indicates that a required
	// authentication mechanism is unsupported.
	ErrRequiredAuthMechanismUnsupported = errors.New("required auth mechanism unsupported")
)

// Login uses the provided client connection and credentials to login to the
// IMAP server using plaintext authentication. Most servers will reject logins
// unless TLS is used.
func Login(c *client.Client, username string, password string, logger zerolog.Logger) error {

	if c == nil {
		errMsg := fmt.Sprintf(
			"invalid (nil) client received while attempting login for account %s",
			username,
		)

		logger.Error().Msg(errMsg)

		return fmt.Errorf(errMsg)
	}

	// Due to logic applied during connection establishment this is highly
	// unlikely to be true, but on the mischance that it is we issue a
	// warning.
	if !c.IsTLS() {
		logger.Warn().Msg("WARNING: Connection to server is insecure (TLS is not enabled)")
	}

	capabilities, err := c.Capability()
	if err != nil {
		logger.Debug().Err(err).Msg("Unable to list server capabilities")
		return fmt.Errorf(
			"unable to list server capabilities: %w",
			err,
		)
	}

	for k, v := range capabilities {
		if v {
			logger.Debug().Msgf("Capability: %v", k)
		}
	}

	// Make sure that LOGINDISABLED capability has not been set. When this is
	// advertised by the server the LOGIN command is rejected. A common
	// scenario where this capability is set is when the client has not yet
	// established a secure TLS connection.
	// https://datatracker.ietf.org/doc/html/rfc3501#section-6.2.3
	loginDisabled, capErr := c.Support(IMAPv4CapabilityLoginDisabled)
	if capErr != nil {
		return fmt.Errorf(
			"failed to detect support for logins: %w",
			capErr,
		)
	}

	if loginDisabled {
		return fmt.Errorf(
			"server has disabled logins: %w",
			client.ErrLoginDisabled,
		)
	}

	logger.Debug().Msg("Logging in")
	if err := c.Login(username, password); err != nil {
		errMsg := "login error occurred"
		logger.Error().Err(err).Msg(errMsg)

		return fmt.Errorf("%s: %w", errMsg, err)
	}
	logger.Debug().Msg("Logged in")

	return nil

}

// OAuth2ClientCredsAuth uses the provided client connection and OAuth2
// settings for Client Credentials flow authentication.
//
// The XOAUTH2 authentication mechanism is used as described in
// https://developers.google.com/gmail/xoauth2_protocol and
// https://learn.microsoft.com/en-us/exchange/client-developer/legacy-protocols/how-to-authenticate-an-imap-pop-smtp-application-by-using-oauth#sasl-xoauth2
func OAuth2ClientCredsAuth(
	ctx context.Context,
	imapClient *client.Client,
	mailbox string,
	clientID string,
	clientSecret string,
	scopes []string,
	tokenEndpointURL string,
	logger zerolog.Logger,
) error {

	if imapClient == nil {
		errMsg := fmt.Sprintf(
			"invalid (nil) client received while attempting login for client ID %s",
			clientID,
		)

		logger.Error().Msg(errMsg)

		return fmt.Errorf(errMsg)
	}

	// Due to logic applied during connection establishment this is highly
	// unlikely to be true, but on the mischance that it is we issue a
	// warning.
	if !imapClient.IsTLS() {
		logger.Warn().Msg("WARNING: Connection to server is insecure (TLS is not enabled)")
	}

	oauth2Config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenEndpointURL,
		Scopes:       scopes,
	}

	logger.Debug().Msg("Acquiring fresh token")
	token, err := oauth2Config.Token(ctx)
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to retrieve token")
		return fmt.Errorf(
			"failed to retrieve token: %w",
			err,
		)
	}
	logger.Debug().
		Str("token_expiration", token.Expiry.Format(time.RFC3339)).
		Str("token_type", token.Type()).
		Msg("Token acquired")

	accessToken := token.AccessToken

	// NOTE: Security concern. Perhaps log at trace level to console instead
	// of the normal logger output path (e.g., if normally a file, send to
	// console to prevent default logging to file?).
	//
	// logger.Debug().
	// 	Str("access_token", accessToken).
	// 	Msg("AccessToken")

	capabilities, err := imapClient.Capability()
	if err != nil {
		logger.Debug().Err(err).Msg("Unable to list server capabilities")
		return fmt.Errorf(
			"unable to list server capabilities: %w",
			err,
		)
	}

	for k, v := range capabilities {
		if v {
			logger.Debug().Msgf("Capability: %v", k)
		}
	}

	// Assert that the expected authentication mechanism is supported.
	xoauth2Supported, err := imapClient.SupportAuth(sasl.Xoauth2)
	if err != nil {
		logger.Debug().
			Err(err).
			Str("mechanism", sasl.Xoauth2).
			Msg("Failed to confirm mechanism support")

		return fmt.Errorf(
			"failed to confirm support for mechanism %s: %w",
			sasl.Xoauth2,
			err,
		)
	}

	if !xoauth2Supported {
		logger.Debug().
			Str("mechanism", sasl.Xoauth2).
			Msg("Server does not support required mechanism")

		return fmt.Errorf(
			"%s auth mechanism unavailable: %w",
			sasl.Xoauth2,
			ErrRequiredAuthMechanismUnsupported,
		)
	}

	logger.Debug().
		Str("mechanism", sasl.Xoauth2).
		Msg("Server supports required mechanism")

	// Login to the IMAP server with XOAUTH2
	saslClient := sasl.NewXoauth2Client(mailbox, accessToken)
	if err := imapClient.Authenticate(saslClient); err != nil {
		logger.Debug().
			Err(err).
			Str("mechanism", sasl.Xoauth2).
			Str("client_id", clientID).
			Str("mailbox", mailbox).
			Msg("Failed to authenticate.")

		return fmt.Errorf(
			"failed to authenticate: %w",
			err,
		)
	}

	logger.Debug().Msg("Logged in")

	connState := int(imapClient.State())
	logger.Debug().Int("connection_state", connState).Msg("Connection state")

	return nil

}
