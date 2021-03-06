// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import "strings"

// textileFormattingCharReplacements is a map of replacement characters that
// are known to interfere with Textile formatted reports generated by the
// list-emails tool.
var textileFormattingCharReplacements = map[rune]string{

	// replace pipe characters with HTML entity equivalent
	'|': "&#124;",

	// replace octothorpe (aka, "hashtag", "pound sign" or "number sign")
	// character with HTML entity equivalent
	'#': "&#35;",
}

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

// ReplaceUsingSubMap accepts an original string and a map of character
// replacements. For every matching character found in the original string a
// replacement string (e.g, an HTML entity value) from the map is used in its
// place. A modified copy of the original string is returned.
func ReplaceUsingSubMap(s string, substMap map[rune]string) string {

	var b strings.Builder

	for _, c := range s {
		sub, hasSub := substMap[c]
		switch {
		case hasSub:
			b.WriteString(sub)
		default:
			b.WriteRune(c)
		}
	}

	return b.String()

}

// ReplaceTextileFormatCharacters accepts an original string and uses the
// ReplaceUsingSubMap function to replace any characters specific to Textile
// formatting. A modified copy of the original string is returned.
func ReplaceTextileFormatCharacters(s string) string {
	return ReplaceUsingSubMap(s, textileFormattingCharReplacements)
}
