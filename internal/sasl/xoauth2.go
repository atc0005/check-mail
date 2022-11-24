// Copyright 2016 emersion
// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package sasl

import (
	"encoding/json"
	"fmt"

	gosasl "github.com/emersion/go-sasl"
)

// Add an "implements assertion" to fail the build if the implementation of
// the upstream Client interface isn't correct.
var _ gosasl.Client = (*xoauth2Client)(nil)

// Xoauth2 is the IMAP authentication mechanism name.
const Xoauth2 = "XOAUTH2"

// Xoauth2Error represents an error encountered during an XOAUTH2
// authentication attempt.
type Xoauth2Error struct {
	Status  string `json:"status"`
	Schemes string `json:"schemes"`
	Scope   string `json:"scope"`
}

// Error implements the stdlib error interface.
func (err *Xoauth2Error) Error() string {
	return fmt.Sprintf("XOAUTH2 authentication error (%v)", err.Status)
}

// xoauth2Client represents a client used to perform challenge-response
// authentication using the XOAUTH2 mechanism.
type xoauth2Client struct {
	Username string
	Token    string
}

// Start begins SASL authentication with the server. It returns the
// authentication mechanism name and "initial response" data consisting of a
// given username and access token encoded in SASL XOAUTH2 format *but*
// without base64 encoding. The base64 encoding step occurs later as part of
// submitting IMAP commands.
func (a *xoauth2Client) Start() (mech string, ir []byte, err error) {
	mech = Xoauth2
	ir = []byte("user=" + a.Username + "\x01auth=Bearer " + a.Token + "\x01\x01")
	return
}

// Next continues the challenge-response authentication. A non-nil error
// causes the client to abort the authentication attempt.
func (a *xoauth2Client) Next(challenge []byte) ([]byte, error) {
	// Server sent an error response
	xoauth2Err := &Xoauth2Error{}
	if err := json.Unmarshal(challenge, xoauth2Err); err != nil {
		return nil, err
	}

	return nil, xoauth2Err
}

// NewXoauth2Client provides an implementation of the XOAUTH2 authentication
// mechanism, as described in
// https://developers.google.com/gmail/xoauth2_protocol and
// https://learn.microsoft.com/en-us/exchange/client-developer/legacy-protocols/how-to-authenticate-an-imap-pop-smtp-application-by-using-oauth#sasl-xoauth2
//
// NOTE: The required base64 encoding of the XOAUTH2 string is performed by
// the emersion/go-imap library as part of submitting the AUTHENTICATE
// command.
func NewXoauth2Client(username, token string) Client {
	return &xoauth2Client{username, token}
}
