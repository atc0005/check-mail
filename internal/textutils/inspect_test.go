// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import (
	"errors"
	"os"
	"testing"
)

// TestInspectAstralUnicodeStrings tests the InspectString function to assert
// that no errors are encountered as part of generating a summary table from
// test input strings (which may contain characters outside of the UTF8MB3
// character set).
func TestInspectAstralUnicodeStrings(t *testing.T) {

	for _, v := range unicodeAstralTestStrings {

		var want error
		got := InspectString(v.original, os.Stderr)

		if !errors.Is(got, want) {
			t.Error("Expected", want, "Got", got)
		}
	}
}

// TestInspectFormattingStrings tests the InspectString function to assert
// that no errors are encountered as part of generating a summary table from
// test input strings specific to Textile formatting characters. These test
// strings are not intended to specifically test handling of characters
// outside of the UTF8MB3 character set, instead leaving that to another more
// focused test case.
func TestInspectFormattingStrings(t *testing.T) {

	for _, v := range formattingTestStrings {

		var want error
		got := InspectString(v.original, os.Stderr)

		if !errors.Is(got, want) {
			t.Error("Expected", want, "Got", got)
		}
	}
}
