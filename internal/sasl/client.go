// Copyright 2016 emersion
// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package sasl

// Client interface to perform challenge-response authentication.
type Client interface {
	// Begins SASL authentication with the server. It returns the
	// authentication mechanism name and "initial response" data (if required by
	// the selected mechanism). A non-nil error causes the client to abort the
	// authentication attempt.
	//
	// A nil ir value is different from a zero-length value. The nil value
	// indicates that the selected mechanism does not use an initial response,
	// while a zero-length value indicates an empty initial response, which must
	// be sent to the server.
	Start() (mech string, ir []byte, err error)

	// Continues challenge-response authentication. A non-nil error causes
	// the client to abort the authentication attempt.
	Next(challenge []byte) (response []byte, err error)
}
