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
func setSummary(accounts []config.MailAccount, nes *nagios.Plugin) {

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
		"%s: No messages found in specified folders for accounts: %v",
		nagios.StateOKLabel,
		accounts,
	)

	for _, account := range accounts {
		// Building with `go build -gcflags=all=-d=loopvar=2` identified this
		// loop as compiling differently with Go 1.22 (per-iteration) loop
		// semantics.
		//
		// As a workaround, we create a new variable for each iteration to
		// work around potential issues with Go versions prior to Go 1.22.
		//
		// NOTE: Not needed as of Go 1.22.
		//
		// account := account

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
