// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	zlog "github.com/rs/zerolog/log"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/check-mail/internal/mbxs"
	"github.com/atc0005/check-mail/internal/textutils"
	"github.com/atc0005/go-nagios"
)

// A baseline starting point for allocating slices to match "about" how many
// mailboxes will need to be checked on a remote mail server. This is defined
// here instead of hard-coding slice preallocation values.
const mailboxCountGuesstimate int = 30

func main() {

	// Start off assuming all is well, adjust as we go.
	var nagiosExitState = nagios.ExitState{
		LastError:      nil,
		ExitStatusCode: nagios.StateOKExitCode,
	}

	// defer this from the start so it is the last deferred function to run
	defer nagiosExitState.ReturnCheckResults()

	// Setup configuration by parsing user-provided flags
	cfg, cfgErr := config.New()
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())

		return

	case cfgErr != nil:
		// We're using the standalone Err function from rs/zerolog/log as we
		// do not have a working configuration.
		zlog.Err(cfgErr).Msg("Error initializing application")
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error initializing application",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.LastError = cfgErr
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	if cfg.EmitBranding {
		// If enabled, show application details at end of notification
		nagiosExitState.BrandingCallback = config.Branding("Notification generated by ")
	}

	server := fmt.Sprintf("%s:%d", cfg.Server, cfg.Port)

	cfg.Log.Debug().Msg("connecting to remote server")
	c, err := client.DialTLS(server, nil)
	if err != nil {
		cfg.Log.Error().Err(err).Msgf("error connecting to server")
		nagiosExitState.LastError = err
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error connecting to %s",
			nagios.StateCRITICALLabel,
			server,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode
		return
	}
	cfg.Log.Debug().Msg("Connected")

	cfg.Log.Debug().Msg("Logging in")
	if err := c.Login(cfg.Username, cfg.Password); err != nil {
		cfg.Log.Error().Err(err).Msg("Login error occurred")
		nagiosExitState.LastError = err
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Login error occurred",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode
		return
	}
	cfg.Log.Debug().Msg("Logged in")

	cfg.Log.Debug().Msg("Defer logout")
	// At this point in the code we are connected to the remote server and are
	// also logged in with a valid account. Calling os.Exit(X) at this point
	// will cause any deferred functions to be skipped, so we instead use
	// nagiosExitState.ReturnCheckResults() anywhere we would have called
	// os.Exit() with the intended status code. This allows us to safely defer
	// a Logout call here and have a reasonable expectation that it will both
	// run AND that we'll have an opportunity to report those logout issues as
	// this application exits.
	defer func(accountName string) {
		cfg.Log.Debug().Msgf("%s: Logging out", accountName)
		if err := c.Logout(); err != nil {
			cfg.Log.Error().Err(err).Msgf("%s: Failed to log out", accountName)
			nagiosExitState.LastError = err
			nagiosExitState.ServiceOutput = fmt.Sprintf(
				"%s: Error logging out",
				nagios.StateWARNINGLabel,
			)
			nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode
			return
		}
		cfg.Log.Debug().Msgf("%s: Logged out", accountName)
	}(cfg.Username)

	// Generate background job to list mailboxes, send down channel until done
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	// NOTE: This goroutine shuts down once c.List() finishes its work
	go func() {
		cfg.Log.Debug().Msg("Running c.List() to fetch a list of available mailboxes")
		done <- c.List("", "*", mailboxes)
	}()

	var mailboxesList = make([]string, 0, mailboxCountGuesstimate)
	for m := range mailboxes {
		cfg.Log.Debug().Msg("collected mailbox from channel")
		mailboxesList = append(mailboxesList, m.Name)
	}

	if err := <-done; err != nil {
		cfg.Log.Error().Err(err).Msg("Error occurred listing mailboxes")
		nagiosExitState.LastError = err
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error occurred listing mailboxes",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode
		return
	}

	cfg.Log.Debug().Msg("no errors encountered listing mailboxes")

	// Prove that our slice is intact
	for _, m := range mailboxesList {
		cfg.Log.Debug().Str("mailbox", m).Msg("")
	}

	// Confirm that requested folders are present on server
	var validatedMailboxesList = make([]string, 0, mailboxCountGuesstimate)
	for _, folder := range cfg.Folders {
		cfg.Log.Debug().Str("mailbox", folder).Msg("Processing requested folder")

		// At this point we are looping over the requested folders, but
		// haven't yet confirmed that they exist as mailboxes on the remote
		// server.

		if strings.ToLower(folder) == "inbox" {

			// NOTE: The "inbox" mailbox/folder name is NOT case-sensitive,
			// but *all* others should be considered case-sensitive. We should
			// be able to safely skip validating "inbox" here since it is a
			// required mailbox/folder name, but all the same we will play it
			// safe and perform a case-insensitive check for a match.
			cfg.Log.Debug().Str("mailbox", folder).Msg("Performing case-insensitive validation")
			if textutils.InList(folder, mailboxesList, true) {
				validatedMailboxesList = append(validatedMailboxesList, folder)
			}

			continue
		}

		cfg.Log.Debug().Str("mailbox", folder).Msg("Performing case-sensitive validation")
		if !textutils.InList(folder, mailboxesList, false) {
			cfg.Log.Error().Str("mailbox", folder).Bool("found", false).Msg("")
			nagiosExitState.LastError = fmt.Errorf("mailbox not found: %q", folder)
			nagiosExitState.ServiceOutput = fmt.Sprintf(
				"%s: Mailbox not found",
				nagios.StateCRITICALLabel,
			)
			nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

			return
		}

		// At this point we have confirmed that the requested folder to
		// monitor is in the list of folders found on the server
		cfg.Log.Debug().Str("mailbox", folder).Bool("found", true).Msg("")
		validatedMailboxesList = append(validatedMailboxesList, folder)

	}

	// At this point we have created a list of validated mailboxes. Process
	// them to determine number of emails within each of them. Based on our
	// existing check and manual processing schedule, we normally see
	// somewhere between 1 and 5 mail items for normal accounts and under 30
	// for heavily spammed accounts. Preallocating the results slice with a
	// midrange starting value for now, but keeping the initial length at 0
	// to allow append() to work as expected.
	var results = make(mbxs.MailboxCheckResults, 0, 10)
	for _, folder := range validatedMailboxesList {

		cfg.Log.Debug().Msg("Selecting mailbox")
		mailbox, err := c.Select(folder, false)
		if err != nil {
			cfg.Log.Error().Err(err).Str("mailbox", folder).Msg("Error occurred selecting mailbox")
			nagiosExitState.LastError = err
			nagiosExitState.ServiceOutput = fmt.Sprintf(
				"%s: Error occurred selecting mailbox %s",
				nagios.StateCRITICALLabel,
				folder,
			)
			nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode
			return
		}

		cfg.Log.Debug().Str("mailbox", folder).Msgf("Mailbox flags: %v", mailbox.Flags)

		cfg.Log.Debug().Msgf("%d mail items found in %s for %s",
			mailbox.Messages, folder, cfg.Username)

		results = append(results, mbxs.MailboxCheckResult{
			MailboxName: folder,
			ItemsFound:  int(mailbox.Messages),
		})
	}

	// Evaluate whether anything was found and sound an alert if so
	if results.GotMail() {
		cfg.Log.Debug().Msgf("%d messages found: %s",
			results.TotalMessagesFound(),
			results.MessagesFoundSummary(),
		)
		nagiosExitState.LastError = nil
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: %s: %d messages found: %s",
			nagios.StateWARNINGLabel,
			cfg.Username,
			results.TotalMessagesFound(),
			results.MessagesFoundSummary(),
		)
		nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode
		return
	}

	// Give the all clear: no mail was found
	cfg.Log.Debug().Msg("No messages found to report")
	nagiosExitState.LastError = nil
	nagiosExitState.ServiceOutput = fmt.Sprintf(
		"%s: %s: No messages found in folders: %s",
		nagios.StateOKLabel,
		cfg.Username,
		cfg.Folders.String(),
	)
	nagiosExitState.ExitStatusCode = nagios.StateOKExitCode
	// implied return here :)

}
