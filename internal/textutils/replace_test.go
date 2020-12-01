// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import "testing"

func TestReplaceAstralUnicode(t *testing.T) {

	for _, v := range testStrings {

		want := v.modified
		got := ReplaceAstralUnicode(v.original, EmojiScissors)

		if got != want {
			t.Error("Expected", want, "Got", got)
		}
	}

}
