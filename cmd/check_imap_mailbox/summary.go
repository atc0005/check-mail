// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/atc0005/go-nagios"
)

// setSummary customizes nagios.ExitState's ServiceOutput and
// LongServiceOutput based on number of user-specified accounts.
func setSummary(accounts []config.MailAccount, nes *nagios.ExitState) {

	if len(accounts) == 1 {
		nes.ServiceOutput = fmt.Sprintf(
			"%s: %s: No messages found in folders: %s",
			nagios.StateOKLabel,
			accounts[0].Username,
			accounts[0].Folders.String(),
		)

		// We're done here. Not much to say if only checking one account.
		return
	}

	nes.ServiceOutput = fmt.Sprintf(
		"%s: %s: No messages found in specified folders for accounts: %v",
		nagios.StateOKLabel,
		accounts[0].Username,
		accounts,
	)

	for _, account := range accounts {
		accountSummary := fmt.Sprintf(
			"* Account: %s%s** Folders: %s%s%s",
			account.Username,
			nagios.CheckOutputEOL,
			account.Folders.String(),
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)
		nes.LongServiceOutput += accountSummary
	}

}
