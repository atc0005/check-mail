// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/atc0005/check-mail/logging"
	"github.com/atc0005/go-nagios"
)

// NagiosExitState represents the last known execution state of the
// application, including the most recent error and the final intended plugin
// state.
// TODO: Refine this further and consider moving to the atc0005/go-nagios
// package.
type NagiosExitState struct {
	LastError  error
	StatusCode int
	Message    string
}

// A baseline starting point for allocating slices to match "about" how many
// mailboxes will need to be checked on a remote mail server. This is defined
// here instead of hard-coding slice preallocation values.
const mailboxCountGuesstimate int = 30

func main() {

	// Start off assuming all is well, adjust as we go.
	var nagiosExitState = NagiosExitState{
		LastError:  nil,
		StatusCode: nagios.StateOK,
	}

	// Nagios relies on plugin exit codes to determine success/failure of
	// checks. The approach that is most often used with other languages is to
	// use something like Using os.Exit() directly and force an early exit of
	// the application with an explicit exit code. Using os.Exit() directly in
	// Go does not run deferred functions; other Go-based plugins that do not
	// rely on deferring function calls may get away with using os.Exit(), but
	// introducing new dependencies could introduce problems.
	//
	// We attempt to explicitly allow deferred functions to work as intended
	// by queuing up values as the app runs and then have this block of code
	// scheduled (deferred) to process those queued values, including the
	// intended plugin exit state. Since this codeblock runs as the last step
	// in the application, it can safely call os.Exit() to set the desired
	// exit code without blocking other deferred functions from running.
	defer func() {
		fmt.Println(nagiosExitState.Message)
		if nagiosExitState.LastError != nil {
			fmt.Printf("\nAdditional details: %v\n", nagiosExitState.LastError)
		}
		os.Exit(nagiosExitState.StatusCode)
	}()

	config := Config{}

	flag.Var(&config.Folders, "folders", "Folders to check for mail. This value is provided a comma-separated list.")
	flag.StringVar(&config.Username, "username", "", "The account used to login to the remote mail server. This is often in the form of an email address.")
	flag.StringVar(&config.Password, "password", "", "The remote mail server account password.")
	flag.StringVar(&config.Server, "server", "", "The fully-qualified domain name of the remote mail server.")
	flag.IntVar(&config.Port, "port", 993, "TCP port used to connect to the remote server. This is usually the same port used for TLS encrypted IMAP connections.")
	flag.StringVar(&config.LoggingLevel, "log-level", "info", "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace.")

	// parse flag definitions from the argument list
	flag.Parse()

	if err := config.Validate(); err != nil {
		log.Err(err).Msg("Error validating configuration")
		nagiosExitState.Message = "Error validating configuration"
		nagiosExitState.LastError = err
		nagiosExitState.StatusCode = nagios.StateCRITICAL
		return
	}

	// Note: Nagios doesn't look at stderr, only stdout. We have to make sure
	// that only whatever output is meant for consumption is emitted to stdout
	// and whatever is meant for troubleshooting is sent to stderr. To help
	// keep these two goals separate (and because Nagios doesn't really do
	// anything special with JSON output from plugins), we use stdlib fmt
	// package output functions for Nagios via stdout and logging package for
	// troubleshooting via stderr.
	//
	// Also, set common fields here so that we don't have to repeat them
	// explicitly later. This will hopefully help to standardize the log
	// messages to make them easier to search through later when
	// troubleshooting.
	log := zerolog.New(os.Stderr).With().Caller().
		Str("username", config.Username).
		Str("server", config.Server).
		Int("port", config.Port).
		Str("folders_to_check", config.Folders.String()).Logger()

	if err := logging.SetLoggingLevel(config.LoggingLevel); err != nil {
		log.Err(err).Msg("configuring logging level")
		nagiosExitState.LastError = err
		nagiosExitState.Message = "Error configuring logging level"
		nagiosExitState.StatusCode = nagios.StateCRITICAL
		return
	}

	server := fmt.Sprintf("%s:%d", config.Server, config.Port)

	log.Debug().Msg("connecting to remote server")
	c, err := client.DialTLS(server, nil)
	if err != nil {
		log.Error().Err(err).Msgf("error connecting to server")
		nagiosExitState.LastError = err
		nagiosExitState.Message = fmt.Sprintf("Error connecting to %s", server)
		nagiosExitState.StatusCode = nagios.StateCRITICAL
		return
	}
	log.Debug().Msg("Connected")

	log.Debug().Msg("Logging in")
	if err := c.Login(config.Username, config.Password); err != nil {
		log.Error().Err(err).Msg("Login error occurred")
		nagiosExitState.LastError = err
		nagiosExitState.Message = "Login error occurred"
		nagiosExitState.StatusCode = nagios.StateCRITICAL
		return
	}
	log.Debug().Msg("Logged in")

	log.Debug().Msg("Defer logout")
	// At this point in the code we are connected to the remote server and
	// are also logged in with a valid account. Calling os.Exit(X) at this
	// point will cause any deferred functions to be skipped, so we instead
	// track our "plugin" state via a function-level variable defined earlier
	// and then return early from main() anywhere we would have called
	// os.Exit() with the intended status code. This allows us to safely
	// defer a Logout call here and have a reasonable expectation that it will
	// both run AND that we'll have an opportunity to report those logout
	// issues as this application exits.
	defer func(accountName string) {
		log.Debug().Msg("Logging out")
		if err := c.Logout(); err != nil {
			log.Error().Err(err).Msg("")
			nagiosExitState.LastError = err
			nagiosExitState.Message = "Error logging out"
			nagiosExitState.StatusCode = nagios.StateWARNING
		}
	}(config.Username)

	// Generate background job to list mailboxes, send down channel until done
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	// TODO: Aside from app exit, what shuts down this goroutine?
	go func() {
		log.Debug().Msg("Running c.List() to fetch a list of available mailboxes")
		done <- c.List("", "*", mailboxes)
	}()

	var mailboxesList = make([]string, 0, mailboxCountGuesstimate)
	for m := range mailboxes {
		log.Debug().Msg("collected mailbox from channel")
		mailboxesList = append(mailboxesList, m.Name)
	}

	// At this point we are finished with the done channel? At what point
	// should we *move on* after wrapping up our use of the channel? The
	// official README "client" example shows checking the channel results
	// *after* ranging over it, so presumably it doesn't need to be checked
	// upfront? Why is this?
	if err := <-done; err != nil {
		log.Error().Err(err).Msg("Error occurred listing mailboxes")
		nagiosExitState.LastError = err
		nagiosExitState.Message = "Error occurred listing mailboxes"
		nagiosExitState.StatusCode = nagios.StateCRITICAL
		return
	}

	log.Debug().Msg("no errors encountered listing mailboxes")

	// Prove that our slice is intact
	for _, m := range mailboxesList {
		log.Debug().Str("mailbox", m).Msg("")
	}

	// Confirm that requested folders are present on server
	var validatedMailboxesList = make([]string, 0, mailboxCountGuesstimate)
	for _, folder := range config.Folders {
		log.Debug().Str("mailbox", folder).Msg("Processing requested folder")

		// At this point we are looping over the requested folders, but
		// haven't yet confirmed that they exist as mailboxes on the remote
		// server.

		if strings.ToLower(folder) == "inbox" {

			// NOTE: The "inbox" mailbox/folder name is NOT case-sensitive,
			// but *all* others should be considered case-sensitive. We should
			// be able to safely skip validating "inbox" here since it is a
			// required mailbox/folder name, but all the same we will play it
			// safe and perform a case-insensitive check for a match.
			log.Debug().Str("mailbox", folder).Msg("Performing case-insensitive validation")
			if InList(folder, mailboxesList, true) {
				validatedMailboxesList = append(validatedMailboxesList, folder)
			}

			continue
		}

		log.Debug().Str("mailbox", folder).Msg("Performing case-sensitive validation")
		if InList(folder, mailboxesList, false) {

			// At this point we have confirmed that the requested folder to
			// monitor is in the list of folders found on the server
			log.Debug().Str("mailbox", folder).Bool("found", true).Msg("")
			validatedMailboxesList = append(validatedMailboxesList, folder)

		} else {

			log.Error().Str("mailbox", folder).Bool("found", false).Msg("")
			nagiosExitState.LastError = fmt.Errorf("mailbox not found: %q", folder)
			nagiosExitState.Message = "Mailbox not found"
			nagiosExitState.StatusCode = nagios.StateCRITICAL
			return
		}

	}

	// At this point we have created a list of validated mailboxes. Process
	// them to determine number of emails within each of them. Based on our
	// existing check and manual processing schedule, we normally see
	// somewhere between 1 and 5 mail items for normal accounts and under 30
	// for heavily spammed accounts. Preallocating the results slice with a
	// midrange starting value for now, but keeping the initial length at 0
	// to allow append() to work as expected.
	var results = make(mailboxCheckResults, 0, 10)
	for _, folder := range validatedMailboxesList {

		log.Debug().Msg("Selecting mailbox")
		mailbox, err := c.Select(folder, false)
		if err != nil {
			log.Error().Err(err).Str("mailbox", folder).Msg("Error occurred selecting mailbox")
			nagiosExitState.LastError = err
			nagiosExitState.Message = fmt.Sprintf("Error occurred selecting mailbox %s", folder)
			nagiosExitState.StatusCode = nagios.StateCRITICAL
			return
		}

		log.Debug().Str("mailbox", folder).Msgf("Mailbox flags: %v", mailbox.Flags)

		log.Debug().Msgf("%d mail items found in %s for %s",
			mailbox.Messages, folder, config.Username)

		results = append(results, mailboxCheckResult{
			mailboxName: folder,
			itemsFound:  int(mailbox.Messages),
		})
	}

	// Evaluate whether anything was found and sound an alert if so
	if results.GotMail() {
		log.Debug().Msgf("%d messages found: %s",
			results.TotalMessagesFound(),
			results.MessagesFoundSummary(),
		)
		nagiosExitState.LastError = nil
		nagiosExitState.Message = fmt.Sprintf("WARNING: %s: %d messages found: %s",
			config.Username,
			results.TotalMessagesFound(),
			results.MessagesFoundSummary(),
		)
		nagiosExitState.StatusCode = nagios.StateWARNING
		return
	}

	// Give the all clear: no mail was found
	log.Debug().Msg("No messages found to report")
	nagiosExitState.LastError = nil
	nagiosExitState.Message = fmt.Sprintf("OK: %s: No messages found in folders: %q",
		config.Username,
		config.Folders,
	)
	nagiosExitState.StatusCode = nagios.StateOK

}
