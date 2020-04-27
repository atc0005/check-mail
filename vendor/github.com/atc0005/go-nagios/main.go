// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package nagios is a small collection of common types and package-level
// variables intended for use with various plugins to reduce code duplication.
package nagios

// Nagios plugin/service check states. These constants replicate the values
// from utils.sh which is normally found at one of these two locations,
// depending on which Linux distribution you're using:
//
//     /usr/lib/nagios/plugins/utils.sh
//     /usr/local/nagios/libexec/utils.sh
//
// See also http://nagios-plugins.org/doc/guidelines.html
const (
	StateOK        int = 0
	StateWARNING   int = 1
	StateCRITICAL  int = 2
	StateUNKNOWN   int = 3
	StateDEPENDENT int = 4
)
