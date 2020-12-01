// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import (
	"os"
	"testing"
)

func TestInspectString(t *testing.T) {

	for _, v := range testStrings {

		var want error = nil
		got := InspectString(v.original, os.Stderr)

		if got != want {
			t.Error("Expected", want, "Got", got)
		}
	}
}
