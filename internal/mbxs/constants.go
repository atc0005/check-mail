// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import "github.com/atc0005/check-mail/internal/textutils"

// DefaultReplacementString is used when replacing Unicode characters
// incompatible with a target character set. The common use case is
// substituting Unicode characters incompatible with the utf8mb3 character
// set.
const DefaultReplacementString string = textutils.EmojiScissors

// Known, named networks used for IMAP connections. These names match the
// network names used by the `net` standard library package.
const (

	// NetTypeTCPAuto indicates that either of IPv4 or IPv6 will be used to
	// establish a connection depending on the specified IP Address.
	NetTypeTCPAuto string = "tcp"

	// NetTypeTCP4 indicates an IPv4-only network.
	NetTypeTCP4 string = "tcp4"

	// NetTypeTCP6 indicates an IPv6-only network.
	NetTypeTCP6 string = "tcp6"
)

// IMAP commands
const (

	// The LOGIN command identifies the client to the server and carries the
	// plaintext password authenticating this user.
	// https://datatracker.ietf.org/doc/html/rfc3501#section-6.2.3
	IMAPv4CommandLogin string = "LOGIN"
)

// IMAP capabilities
const (

	// The LOGINDISABLED capability indicates that the LOGIN command is not
	// permitted. This is most often due to TLS not active for the current
	// connection.
	// https://datatracker.ietf.org/doc/html/rfc3501#section-6.1.1
	IMAPv4CapabilityLoginDisabled string = "LOGINDISABLED"
)
