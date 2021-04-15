// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import (
	"fmt"
	"strings"
	"time"

	"github.com/atc0005/check-mail/internal/textutils"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog"
)

// ListMailboxes lists mailboxes associated with the logged in user account
// (by way of an IMAP client connection).
func ListMailboxes(c *client.Client, logger zerolog.Logger) ([]string, error) {

	// Generate background job to list mailboxes, send down channel until done
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	// NOTE: This goroutine shuts down once c.List() finishes its work
	go func() {
		logger.Debug().Msg("Running c.List() to fetch a list of available mailboxes")
		done <- c.List("", "*", mailboxes)
	}()

	var mailboxesList = make([]string, 0, mailboxCountGuesstimate)
	for m := range mailboxes {
		logger.Debug().Msg("collected mailbox from channel")
		mailboxesList = append(mailboxesList, m.Name)
	}

	if err := <-done; err != nil {
		logger.Error().Err(err).Msg("Error occurred listing mailboxes")

		return nil, err
	}

	logger.Debug().Msg("no errors encountered listing mailboxes")
	return mailboxesList, nil

}

// ValidateMailboxesList receives a list of requested mailboxes and returns a
// list of mailboxes from that list which have been confirmed to be present
// for the associated user account.
func ValidateMailboxesList(c *client.Client, userMBXList []string, logger zerolog.Logger) ([]string, error) {

	// Get list of mailboxes on server to compare against user mailbox list.
	serverMBXList, listErr := ListMailboxes(c, logger)
	if listErr != nil {
		return nil, listErr
	}

	// List out detected mailboxes for debugging purposes.
	for _, m := range serverMBXList {
		logger.Debug().Str("mailbox", m).Msg("")
	}

	// Confirm that requested folders are present on server
	validatedMBXList := make([]string, 0, len(userMBXList))

	for _, mbx := range userMBXList {
		logger.Debug().Str("mailbox", mbx).Msg("Processing requested folder")

		// At this point we are looping over the user requested
		// folders/mailboxes, but haven't yet confirmed that they exist as
		// mailboxes on the remote server.

		if strings.ToLower(mbx) == "inbox" ||
			strings.HasPrefix(strings.ToLower(mbx), "inbox/") {

			// NOTE: The "inbox" mailbox/folder name is NOT case-sensitive,
			// but *all* others should be considered case-sensitive. We should
			// be able to safely skip validating "inbox" here since it is a
			// required mailbox/folder name, but all the same we will play it
			// safe and perform a case-insensitive check for a match.
			logger.Debug().Str("mailbox", mbx).Msg("Performing case-insensitive validation")
			if textutils.InList(mbx, serverMBXList, true) {
				validatedMBXList = append(validatedMBXList, mbx)
			}

			continue
		}

		logger.Debug().Str("mailbox", mbx).Msg("Performing case-sensitive validation")
		if !textutils.InList(mbx, serverMBXList, false) {
			logger.Error().Str("mailbox", mbx).Bool("found", false).Msg("")

			return nil, fmt.Errorf("mailbox not found: %q", mbx)
		}

		// At this point we have confirmed that the requested folder to
		// evaluate is in the list of folders found on the server
		logger.Debug().Str("mailbox", mbx).Bool("found", true).Msg("")
		validatedMBXList = append(validatedMBXList, mbx)

	}

	return validatedMBXList, nil

}

// CheckMail generates a listing of emails within the provided (and validated)
// mailbox list for the associated username.
func CheckMail(c *client.Client, username string, validatedMBXList []string, logger zerolog.Logger) (MailboxCheckResults, error) {

	// Process validated mailboxes list to determine number of emails within
	// each of them. Based on our existing check and manual processing
	// schedule, we normally see somewhere between 1 and 5 mail items for
	// normal accounts and under 30 for heavily spammed accounts.
	// Preallocating the results slice with a midrange starting value for now,
	// but keeping the initial length at 0 to allow append() to work as
	// expected.
	results := make(MailboxCheckResults, 0, 10)
	for _, folder := range validatedMBXList {

		logger.Debug().Str("mailbox", folder).Msg("Selecting mailbox")
		mailbox, selectErr := c.Select(folder, false)
		if selectErr != nil {
			logger.Error().
				Err(selectErr).
				Str("mailbox", folder).
				Msg("Error occurred selecting mailbox")

			return nil, fmt.Errorf(
				"%s: error occurred selecting mailbox: %w",
				username,
				selectErr,
			)
		}

		logger.Debug().Str("mailbox", folder).Msgf(
			"Mailbox flags for %q: %v",
			folder,
			mailbox.Flags,
		)

		logger.Info().Msgf("%d mail items found in %q for %s",
			mailbox.Messages, folder, username)

		// List all email messages, if there are any
		if mailbox.Messages == 0 {
			logger.Debug().
				Str("mailbox", folder).
				Uint32("messages_found", mailbox.Messages).
				Msg("no messages found")

			// record "no results" so we can explicitly note this later
			results = append(results, MailboxCheckResult{
				MailboxName: folder,
				ItemsFound:  0,
			})
			continue
		}

		// specify message retrieval range: from 1 to total available
		seqset := new(imap.SeqSet)
		from := uint32(1)
		to := mailbox.Messages
		seqset.AddRange(from, to)

		// room for 10 messages at once
		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			// room for one error response
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		// fmt.Printf("\n\nEmail messages in mailbox %q: \n\n", folder)

		// collect all emails found within mailbox
		messagesFound := make([]Message, 0, 10)
		for msg := range messages {

			var subject string
			switch {
			case !textutils.WithinUTF8MB3Range(msg.Envelope.Subject):
				logger.Debug().Msg("Replacing Astral Unicode characters")
				subject = textutils.ReplaceAstralUnicode(
					msg.Envelope.Subject, DefaultReplacementString)
			default:
				logger.Debug().Msg("Using original subject line")
				subject = msg.Envelope.Subject
			}
			// fmt.Println("*", subject)

			// Replace any output formatting characters that may be present.
			logger.Debug().Msg("Replacing any Textile characters known to cause issues")
			subject = textutils.ReplaceTextileFormatCharacters(subject)

			msgSummary := Message{
				MessageID:             msg.Envelope.MessageId,
				EnvelopeDate:          msg.Envelope.Date,
				EnvelopeDateFormatted: msg.Envelope.Date.Format(time.RFC3339),
				OriginalSubject:       msg.Envelope.Subject,
			}

			// we only set the ModifiedSubject field if the subject line
			// is actually modified from the original value.
			if subject != msg.Envelope.Subject {
				msgSummary.ModifiedSubject = subject
			}

			messagesFound = append(messagesFound, msgSummary)

		}

		// block until we get a response
		if err := <-done; err != nil {
			logger.Error().
				Err(err).
				Str("mailbox", folder).
				Msg("Error occurred listing emails in mailbox")

			return nil, fmt.Errorf(
				"%s: error occurred listing emails in mailbox %s",
				username,
				folder,
			)

		}

		results = append(results, MailboxCheckResult{
			MailboxName: folder,
			ItemsFound:  int(mailbox.Messages),
			Messages:    messagesFound,
		})

	}

	return results, nil

}
