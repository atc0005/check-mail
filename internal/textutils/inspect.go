// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import (
	"fmt"
	"io"
	"text/tabwriter"
	"unicode"
)

// inspectString is a shared helper function that generates a summary table
// from a provided string to help identify Unicode characters incompatible
// with older database character sets (e.g., UTF8MB3). This summary table is
// written to the provided io.Writer interface.
func inspectString(s string, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 8, 8, 4, '\t', 0)

	var status string

	for i, c := range s {

		status = "\xE2\x9B\x94 (no)"
		if c <= UTF8MB3RangeEndRune {
			// status = "\xF0\x9F\x8C\x9E (yes)"
			status = "\xE2\x9C\x85 (yes)"
		}

		// fmt.Printf(
		_, _ = fmt.Fprintf(
			tw,
			"char %d: %c\t"+
				"Decimal: %d\t"+
				"IsSymbol: %t\t"+
				"UTF8MB3 safe: %v\t"+
				"code point: %U\t"+
				"rune literal: %+q\t"+

				// literal bytes in hex format roughly equivalent to what
				// MySQL/MariaDB uses in their error messages.
				// MariaDB [testing]> insert into unicode values ("Win a golden ticket to WooConf in Seattle😍");
				// ERROR 1366 (22007): Incorrect string value: '\xF0\x9F\x98\x8D' for column `testing`.`unicode`.`string` at row 1
				"Hex: % X\n",
			i,
			c,
			c,
			unicode.IsSymbol(c),
			// c <= UTF8MB3RangeEndRune,
			status,
			c,
			c,
			// convert rune to string, then to byte slice
			[]byte(string(c)),
		)
		// }

	}

	_, _ = fmt.Fprintln(w)
	if err := tw.Flush(); err != nil {
		return fmt.Errorf(
			"error occurred flushing tabwriter: %w",
			err,
		)
	}

	return nil

}

// InspectStrings generates a summary table from a provided slice of strings
// to help identify Unicode characters incompatible with older database
// character sets (e.g., UTF8MB3). This summary table is written to the
// provided io.Writer interface.
func InspectStrings(ss []string, w io.Writer) error {

	for i, s := range ss {

		_, _ = fmt.Fprintf(w, "\nstring %d: %q\n", i, s)

		err := inspectString(s, w)
		if err != nil {
			return err
		}

		fmt.Printf("\n\n**************************************************\n\n")
	}

	return nil
}

// InspectString generates a summary table from a provided string to help
// identify Unicode characters incompatible with older database character sets
// (e.g., UTF8MB3). This summary table is written to the provided io.Writer
// interface.
func InspectString(s string, w io.Writer) error {

	_, _ = fmt.Fprintf(w, "\nstring: %q\n", s)

	err := inspectString(s, w)
	if err != nil {
		return err
	}

	fmt.Printf("\n\n**************************************************\n\n")

	return nil
}

// CharsWithinRange indicates whether a provided string contains any
// characters outside of the provided character set range.
func CharsWithinRange(s string, start rune, end rune) bool {
	for _, c := range s {
		if c > end || c < start {
			return false
		}
	}

	return true
}

// WithinUTF8MB3Range indicates whether a provided string contains any
// characters outside of the UTF8MB3 character set range.
func WithinUTF8MB3Range(s string) bool {
	return CharsWithinRange(s, UTF8MB3RangeStartRune, UTF8MB3RangeEndRune)
}
