// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import (
	"net"
)

// Dialer is an implementation of the go-imap/client.Dialer interface. This
// implementation is used to override the default "auto" network type
// selection behavior applied by `Dial*` functions within the go-imap/client
// package.
type Dialer struct {
	// If specified, this value overrides the value supplied to
	// `go-imap/client.Dial*` functions and is used in its place to establish
	// a connection on the user-specified network type.
	NetworkTypeUserOverride string

	// Original value supplied to `go-imap/client.Dial*` functions. Unless
	// overridden by the user, this value is used to establish a connection on
	// the supplied network type.
	NetworkTypeOriginalValue string
}

// Dial implements the go-imap/client.Dialer interface to override the default
// "auto" network type selection behavior applied by `Dial*` functions within
// the go-imap/client package.
func (d *Dialer) Dial(network string, addr string) (net.Conn, error) {

	// Record the requested network type for potential later use
	d.NetworkTypeOriginalValue = network

	if d.NetworkTypeUserOverride != "" {
		network = d.NetworkTypeUserOverride
	}

	return net.Dial(network, addr)

}
