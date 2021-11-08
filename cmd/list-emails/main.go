// Copyright 2020 Adam Chalkley
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
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/check-mail/internal/files"
	"github.com/atc0005/check-mail/internal/mbxs"
)

type exitStatus struct {
	Code int
	Err  error
}

func main() {

	var appExitStatus exitStatus

	defer func(appExitStatus *exitStatus) {
		if appExitStatus.Err != nil {
			fmt.Println(appExitStatus.Err.Error())
		}
		os.Exit(appExitStatus.Code)
	}(&appExitStatus)

	// Setup configuration by parsing user-provided flags
	useConfigFile := true
	useLogFile := true
	cfg, cfgErr := config.New(useConfigFile, useLogFile)
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())
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
	for _, account := range cfg.Accounts {

		fmt.Println("Checking account:", account.Name)

		logger := cfg.Log.With().
			Str("username", account.Username).
			Str("server", account.Server).
			Int("port", account.Port).
			Str("folders_to_check", account.Folders.String()).
			Logger()

		c, connectErr := mbxs.Connect(account.Server, account.Port, cfg.NetworkType, cfg.MinTLSVersion(), logger)
		if connectErr != nil {
			logger.Error().Err(connectErr).Msg("failed to connect to server")

			appExitStatus.Err = connectErr
			appExitStatus.Code = 1

			return
		}

		if loginErr := mbxs.Login(c, account.Username, account.Password, logger); loginErr != nil {
			logger.Error().Err(loginErr).Msg("failed to login to server")

			appExitStatus.Err = loginErr
			appExitStatus.Code = 1

			return
		}

		logger.Debug().Msg("Defer logout")
		// At this point we are connected to the remote server and are also
		// logged in with a valid account. We defer a Logout call here with a
		// reasonable expectation that it will both run AND that we'll have an
		// opportunity to report those logout issues as this application
		// exits. We do not set appExitErr as we do not want to chance
		// overwriting an existing error message returned from a step prior to
		// the final deferred logout attempt.
		defer func(accountName string) {
			logger.Info().Msgf("%s: Logging out", accountName)
			if err := c.Logout(); err != nil {
				logger.Error().Err(err).Msgf("%s: Failed to log out", accountName)
				appExitStatus.Code = 1

				return
			}
			logger.Info().Msgf("%s: Logged out", accountName)
		}(account.Username)

		// Confirm that requested folders are present on server
		validatedMBXList, validateErr := mbxs.ValidateMailboxesList(
			c, account.Folders, logger)
		if validateErr != nil {
			logger.Error().Err(validateErr).Msg("failed to validate mailboxes list")

			appExitStatus.Err = validateErr
			appExitStatus.Code = 1

			return

		}

		results, chkMailErr := mbxs.CheckMail(c, account.Name, validatedMBXList, logger)
		if chkMailErr != nil {
			logger.Error().Err(chkMailErr).Msg("failed to check mail in mailboxes")

			appExitStatus.Err = chkMailErr
			appExitStatus.Code = 1

			return
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

			appExitStatus.Err = fmt.Errorf(
				"failed to generate report: %w",
				reportGenErr,
			)
			appExitStatus.Code = 1

			return
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

	}

	accountsList := make([]string, 0, len(cfg.Accounts))
	for _, account := range cfg.Accounts {
		accountsList = append(accountsList, account.Username)
	}

	fmt.Printf(
		"OK: Successfully generated reports for accounts: %s\n",
		strings.Join(accountsList, ", "),
	)

	// implied return here :)

}
