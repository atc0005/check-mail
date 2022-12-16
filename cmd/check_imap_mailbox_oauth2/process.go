// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"context"
	"fmt"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/check-mail/internal/mbxs"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

func processAccount(
	ctx context.Context,
	account config.MailAccount,
	cfg *config.Config,
	state *nagios.Plugin,
	logger zerolog.Logger,
) (mbxs.MailboxCheckResults, error) {

	c, connectErr := mbxs.Connect(account.Server, account.Port, cfg.NetworkType, cfg.MinTLSVersion(), logger)
	if connectErr != nil {
		logger.Error().Err(connectErr).Msg("error connecting to server")
		state.AddError(connectErr)
		state.ServiceOutput = fmt.Sprintf(
			"%s: Error connecting to %s",
			nagios.StateCRITICALLabel,
			account.Server,
		)
		state.ExitStatusCode = nagios.StateCRITICALExitCode

		return nil, connectErr
	}
	logger.Debug().Msg("Connection established to server")

	// Enable client network command/response logging if global logging
	// level indicates user wishes to see verbose details.
	if zerolog.GlobalLevel() == zerolog.DebugLevel ||
		zerolog.GlobalLevel() == zerolog.TraceLevel {
		c.SetDebug(&logger)
	}

	// https://github.com/emersion/go-imap#client-
	logger.Debug().Msg("Defer closing connection to server")
	defer func() {
		logger.Debug().Msg("Calling Logout to gracefully close connection to server")
		if err := c.Logout(); err != nil {
			logger.Error().Err(err).Msg("failed to close connection to server")
		}
		logger.Debug().Msg("Connection to server successfully closed")
	}()

	loginErr := mbxs.OAuth2ClientCredsAuth(
		ctx,
		c,
		account.OAuth2Settings.SharedMailbox,
		account.OAuth2Settings.ClientID,
		account.OAuth2Settings.ClientSecret,
		account.OAuth2Settings.Scopes,
		account.OAuth2Settings.TokenURL,
		cfg.RetrievalAttempts(),
		logger,
	)
	if loginErr != nil {
		logger.Error().Err(loginErr).Msg("Login error occurred")
		state.AddError(loginErr)
		state.ServiceOutput = fmt.Sprintf(
			"%s: Login error occurred",
			nagios.StateCRITICALLabel,
		)
		state.ExitStatusCode = nagios.StateCRITICALExitCode

		return nil, loginErr
	}
	logger.Debug().Msg("Successfully logged in")

	// Confirm that requested folders are present on server
	validatedMBXList, validateErr := mbxs.ValidateMailboxesList(
		c, account.Folders, logger)
	if validateErr != nil {
		state.AddError(validateErr)
		state.ServiceOutput = fmt.Sprintf(
			"%s: %s",
			nagios.StateCRITICALLabel,
			validateErr.Error(),
		)
		state.ExitStatusCode = nagios.StateCRITICALExitCode

		return nil, validateErr

	}

	results, chkMailErr := mbxs.CheckMail(c, account.OAuth2Settings.SharedMailbox, validatedMBXList, logger)
	if chkMailErr != nil {
		state.AddError(chkMailErr)
		state.ServiceOutput = fmt.Sprintf(
			"%s: Error occurred checking mail: %s",
			nagios.StateCRITICALLabel,
			chkMailErr.Error(),
		)
		state.ExitStatusCode = nagios.StateCRITICALExitCode

		return nil, chkMailErr
	}

	return results, nil

}
