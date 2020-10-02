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
)

// MailboxCheckResult notes how many mail items were found for a specific
// mailbox
type MailboxCheckResult struct {
	MailboxName string
	ItemsFound  int
}

// MailboxCheckResults represents a collection of all results from mailbox
// checks.
type MailboxCheckResults []MailboxCheckResult

// GotMail returns true if mail was found in checked mailboxes or false if not.
func (mcr MailboxCheckResults) GotMail() bool {
	for _, result := range mcr {
		if result.ItemsFound > 0 {
			return true
		}
	}
	return false
}

// TotalMessagesFound returns a count of all messages found across all checked
// mailboxes.
func (mcr MailboxCheckResults) TotalMessagesFound() int {
	var total int
	for _, result := range mcr {
		total += result.ItemsFound
	}
	return total
}

// MessagesFoundSummary returns a one-line summary of the mail items found in
// checked mailboxes.
func (mcr MailboxCheckResults) MessagesFoundSummary() string {
	var summary string
	for index, result := range mcr {
		summary += fmt.Sprintf("%s(%d)", result.MailboxName, result.ItemsFound)
		if index < (len(mcr) - 1) {
			// Append separator chars if not processing the last index item
			summary += ", "
		}
	}
	return summary
}

// InList is a helper function to emulate Python's `if "x"
// in list:` functionality. The caller can optionally ignore case of compared
// items.
func InList(needle string, haystack []string, ignoreCase bool) bool {
	for _, item := range haystack {

		if ignoreCase {
			item = strings.ToLower(item)
			needle = strings.ToLower(needle)
		}

		if item == needle {
			return true
		}
	}
	return false
}
