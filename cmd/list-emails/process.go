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
	"strings"
	"time"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/check-mail/internal/files"
	"github.com/atc0005/check-mail/internal/mbxs"
	"github.com/rs/zerolog"
)

func processAccount(
	ctx context.Context,
	account config.MailAccount,
	cfg *config.Config,
	logger zerolog.Logger,
) error {

	switch account.AuthType {
	case config.AuthTypeBasic:

		// Add Basic Auth related fields to logger.
		logger = logger.With().
			Str("username", account.Username).
			Logger()

	case config.AuthTypeOAuth2ClientCreds:
		// Add OAuth2 related fields to logger.
		logger = logger.With().
			Str("client_id", account.OAuth2Settings.ClientID).
			Str("scopes", func() string {
				return strings.Join(account.OAuth2Settings.Scopes, ", ")
			}()).
			Str("shared_mailbox", account.OAuth2Settings.SharedMailbox).
			Logger()

	default:
		// This should not occur; config validation should prevent this
		// scenario.
		logger.Error().Msg("invalid authentication type for account")
		return config.ErrInvalidAuthType

	}

	c, connectErr := mbxs.Connect(account.Server, account.Port, cfg.NetworkType, cfg.MinTLSVersion(), logger)
	if connectErr != nil {
		logger.Error().Err(connectErr).Msg("failed to connect to server")
		return connectErr
	}
	logger.Info().Msg("Connection established to server")

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
		logger.Info().Msg("Connection to server successfully closed")
	}()

	switch account.AuthType {

	case config.AuthTypeBasic:
		if loginErr := mbxs.Login(c, account.Username, account.Password, logger); loginErr != nil {
			logger.Error().Err(loginErr).Msg("failed to login to server")
			return loginErr
		}

	case config.AuthTypeOAuth2ClientCreds:
		loginErr := mbxs.OAuth2ClientCredsAuth(
			ctx,
			c,
			account.OAuth2Settings.SharedMailbox,
			account.OAuth2Settings.ClientID,
			account.OAuth2Settings.ClientSecret,
			account.OAuth2Settings.Scopes,
			account.OAuth2Settings.TokenURL,
			logger,
		)
		if loginErr != nil {
			logger.Error().Err(loginErr).Msg("failed to login to server")
			return loginErr
		}
		logger.Debug().Msg("Successfully logged in")

	}

	// Confirm that requested folders are present on server
	validatedMBXList, validateErr := mbxs.ValidateMailboxesList(
		c, account.Folders, logger)
	if validateErr != nil {
		logger.Error().Err(validateErr).Msg("failed to validate mailboxes list")
		return validateErr

	}

	results, chkMailErr := mbxs.CheckMail(c, account.Name, validatedMBXList, logger)
	if chkMailErr != nil {
		logger.Error().Err(chkMailErr).Msg("failed to check mail in mailboxes")
		return chkMailErr
	}

	summaryMsg := fmt.Sprintf("%d messages found: %s",
		results.TotalMessagesFound(),
		results.MessagesFoundSummary(),
	)
	logger.Info().Msg(summaryMsg)

	// Generate report for this account
	reportData := files.ReportData{
		AccountName:           account.Name,
		MailboxCheckResults:   results,
		MessagesFoundSummary:  results.MessagesFoundSummary(),
		ReportTime:            time.Now(),
		UnicodeCharSubstitute: mbxs.DefaultReplacementString,
	}

	reportGenErr := files.GenerateReport(reportData, cfg.ReportFileOutputDir, logger)
	if reportGenErr != nil {

		logger.Error().Err(reportGenErr).Msg("failed to generate report")

		return fmt.Errorf(
			"failed to generate report: %w",
			reportGenErr,
		)
	}

	if cfg.LoggingLevel == config.LogLevelDebug {
		for _, mailbox := range results {
			fmt.Printf(
				"\n\nEmail messages in mailbox %q: \n\n",
				mailbox.MailboxName,
			)

			if len(mailbox.Messages) == 0 {
				fmt.Println("* No messages")
				continue
			}

			for _, msg := range mailbox.Messages {
				fmt.Printf(
					// "{\n\tMessageID: %s\n\tEnvelopeDate: %v\n\tEnvelopeLocalDate: %v\n\tReceivedDate: %v\n\tReceivedLocalDate: %v\n\tOriginalSubject: %s\n\tModifiedSubject: %s\n}\n",
					"{\n\tMessageID: %s\n\t"+
						"EnvelopeDate: %v\n\t"+
						"EnvelopeLocalDate: %v\n\t"+
						"OriginalSubject: %s\n\t"+
						"ModifiedSubject: %s\n}\n",
					msg.MessageID,
					msg.EnvelopeDate,
					msg.EnvelopeDate.Local(),
					// msg.ReceivedDate,
					// msg.ReceivedDate.Local(),
					msg.OriginalSubject,
					msg.ModifiedSubject,
				)
			}

		}

	}

	return nil
}
