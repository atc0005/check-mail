// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

// https://yourbasic.org/golang/multiline-string/
// https://golang.org/ref/spec#Rune_literals
// https://golangbyexample.com/understanding-rune-in-golang/
// https://www.geeksforgeeks.org/rune-in-golang/
// http://www.unicode.org/Public/UCA/9.0.0/allkeys.txt
// https://dev.mysql.com/doc/refman/8.0/en/charset-unicode-sets.html
// https://mathiasbynens.be/notes/mysql-UTF8MB4
// https://apps.timwhitlock.info/emoji/tables/unicode
// https://www.utf8-chartable.de/unicode-utf8-table.pl?start=127808&utf8=string-literal

// https://dev.mysql.com/doc/refman/8.0/en/charset-unicode-sets.html
//
// MySQL implements the xxx_unicode_ci collations according to the Unicode
// Collation Algorithm (UCA) described at
// http://www.unicode.org/reports/tr10/. The collation uses the version-4.0.0
// UCA weight keys: http://www.unicode.org/Public/UCA/4.0.0/allkeys-4.0.0.txt.
// The xxx_unicode_ci collations have only partial support for the Unicode
// Collation Algorithm. Some characters are not supported, and combining marks
// are not fully supported. This affects primarily Vietnamese, Yoruba, and
// some smaller languages such as Navajo. A combined character is considered
// different from the same character written with a single unicode character
// in string comparisons, and the two characters are considered to have a
// different length (for example, as returned by the CHAR_LENGTH() function or
// in result set metadata).
//
// Unicode collations based on UCA versions higher than 4.0.0 include the
// version in the collation name. Examples:
//
// * UTF8MB4_unicode_520_ci is based on UCA 5.2.0 weight keys (http://www.unicode.org/Public/UCA/5.2.0/allkeys.txt),
// * UTF8MB4_0900_ai_ci is based on UCA 9.0.0 weight keys (http://www.unicode.org/Public/UCA/9.0.0/allkeys.txt).

// https://mathiasbynens.be/notes/mysql-UTF8MB4
//
// Since astral symbols (whose code points range from U+010000 to U+10FFFF)
// each consist of four bytes in UTF-8, you cannot store them using MySQL's
// utf8 implementation.

// https://mathiasbynens.be/notes/javascript-encoding#bmp
//
// The first plane (xy is 00) is called the Basic Multilingual Plane or BMP.
// It contains the code points from U+0000 to U+FFFF, which are the most
// frequently used characters.
//
// The other sixteen planes (U+010000 → U+10FFFF) are called supplementary
// planes or astral planes. I won't discuss them here; just remember that
// there are BMP characters and non-BMP characters, the latter of which are
// also known as supplementary characters or astral characters.

// Start and end Unicode code points for the UTF8MB3 character set.
// https://en.wikibooks.org/wiki/Unicode/Character_reference/0000-0FFF
// https://en.wikibooks.org/wiki/Unicode/Character_reference/F000-FFFF
const (
	UTF8MB3RangeStartRune rune = '\u0000'
	UTF8MB3RangeStartInt  int  = 0
	UTF8MB3RangeEndRune   rune = '\uFFFF'
	UTF8MB3RangeEndInt    int  = 65535
)

// Start and end Unicode code points for the UTF8MB4 character set.
// https://en.wikibooks.org/wiki/Unicode/Character_reference/10000-10FFF
// https://en.wikibooks.org/wiki/Unicode/Character_reference/F0000-10FFFF
const (
	UTF8MB4RangeStartRune rune = '\U00010000'
	UTF8MB4RangeStartInt  int  = 65536
	UTF8MB4RangeEndRune   rune = '\U0010FFFF'
	UTF8MB4RangeEndInt    int  = 1114111
)

// common Emoji characters used by this project. Not all are UTF8MB3
// compatible.
const (
	EmojiScissors       string = "\xE2\x9C\x82"
	EmojiNoEntry        string = "\xE2\x9B\x94"
	EmojiRecycle        string = "\xE2\x99\xBB"
	EmojiSunWithFace    string = "\xF0\x9F\x8C\x9E"
	EmojiHeavyCheckMark string = "\xE2\x9C\x85"
	EmojiCrossMark      rune   = '\u274C'
	EmojiCheckMark      rune   = '\u2714'
	EmojiOKButton       rune   = '\U0001F197'
)
