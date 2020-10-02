// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import "strings"

// InList is a helper function to emulate Python's `if "x"
// in list:` functionality. The caller can optionally ignore case of compared
// items.
func InList(needle string, haystack []string, ignoreCase bool) bool {
	for _, item := range haystack {

		if ignoreCase {
			item = strings.ToLower(item)
			needle = strings.ToLower(needle)
		}

		if item == needle {
			return true
		}
	}
	return false
}
