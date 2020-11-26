// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import "strings"

// ReplaceAstralUnicode accepts an original string and a replacement string.
// For every Unicode code point found in the original string that is outside
// of the range of UTF8MB3, the replacement string is used in its place. A
// modified copy of the original string is returned.
func ReplaceAstralUnicode(s string, r string) string {

	var b strings.Builder

	for _, c := range s {
		switch {
		case c > UTF8MB3RangeEndRune:
			b.WriteString(r)
		default:
			b.WriteRune(c)
		}
	}

	return b.String()
}
