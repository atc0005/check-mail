// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import "testing"

// TestReplaceUsingSubMap asserts that the expected Textile formatting
// characters replacement behavior functions as intended by directly testing
// the ReplaceUsingSubMap function.
func TestReplaceUsingSubMap(t *testing.T) {

	for _, v := range formattingTestStrings {

		want := v.modified
		got := ReplaceUsingSubMap(v.original, textileFormattingCharReplacements)

		if got != want {
			t.Error("Expected", want, "Got", got)
		}
	}

}

// TestReplaceTextileCharacters asserts that the expected Textile formatting
// characters replacement behavior functions as intended by using the exported
// & simplified ReplaceTextileCharacters function that wraps the
// ReplaceUsingSubMap function.
func TestReplaceTextileCharacters(t *testing.T) {

	for _, v := range formattingTestStrings {

		want := v.modified
		got := ReplaceTextileFormatCharacters(v.original)

		if got != want {
			t.Error("Expected", want, "Got", got)
		}
	}

}

func TestReplaceAstralUnicode(t *testing.T) {

	for _, v := range unicodeAstralTestStrings {

		want := v.modified
		got := ReplaceAstralUnicode(v.original, EmojiScissors)

		if got != want {
			t.Error("Expected", want, "Got", got)
		}
	}

}
