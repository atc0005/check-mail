// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"crypto/tls"
	"strings"
	"time"
)

// MinTLSVersion returns the applicable `tls.VersionTLS*` numeric constant
// corresponding to the user-specified (or default) TLS version string
// keyword.
func (c Config) MinTLSVersion() uint16 {

	// https://golang.org/pkg/crypto/tls/#pkg-constants
	var tlsVersion uint16

	switch strings.ToLower(c.minTLSVersion) {
	case minTLSVersion10:
		tlsVersion = tls.VersionTLS10
	case minTLSVersion11:
		tlsVersion = tls.VersionTLS11
	case minTLSVersion12:
		tlsVersion = tls.VersionTLS12
	case minTLSVersion13:
		tlsVersion = tls.VersionTLS13
	default:
		tlsVersion = tls.VersionTLS12

	}

	return tlsVersion

}

// MinTLSVersionKeyword returns the user-specified (or default) TLS version
// string keyword.
func (c Config) MinTLSVersionKeyword() string {

	return c.minTLSVersion

}

// SupportedAuthTypes returns the complete list of supported authentication
// types used by applications in this project.
func (c Config) SupportedAuthTypes() []string {
	return []string{
		AuthTypeBasic,
		AuthTypeOAuth2ClientCreds,
	}
}

// Server returns the IMAP server value from the first account entry in the
// collection. This value should be the same for all account entries in the
// collection. For the Plugin app type this is a collection of one entry, for
// the Reporter app type the same server value is recorded for each account
// (convenience).
func (c Config) Server() string {
	return c.Accounts[0].Server
}

// Port returns the IMAP server port value from the first account entry in the
// collection. This value should be the same for all account entries in the
// collection. For the Plugin app type this is a collection of one entry, for
// the Reporter app type the same server value is recorded for each account
// (convenience).
func (c Config) Port() int {
	return c.Accounts[0].Port
}

// AuthType returns the authentication type from the first account entry in
// the collection. This value should be the same for all account entries in
// the collection. For the Plugin app type this is a collection of one entry,
// for the Reporter app type the same server value is recorded for each
// account (convenience).
func (c Config) AuthType() string {
	return c.Accounts[0].AuthType
}

// AccountProcessDelay returns the configured delay between processing
// accounts in a collection.
func (c Config) AccountProcessDelay() time.Duration {

	// TODO: Provide a flag / config file setting for this.
	return defaultAccountProcessDelay
}

// AccountNames returns the collection of names associated with accounts. If
// basic auth is used this will be usernames, if oauth2 is used this will be
// mailbox names. If an account or mailbox is specified as "user@example.com"
// only the "user" portion of that value is returned.
func (c Config) AccountNames() []string {
	accountsList := make([]string, 0, len(c.Accounts))
	for _, account := range c.Accounts {
		switch account.AuthType {
		case AuthTypeBasic:
			name := strings.Split(account.Username, "@")[0]
			accountsList = append(accountsList, name)
		case AuthTypeOAuth2ClientCreds:
			name := strings.Split(account.OAuth2Settings.SharedMailbox, "@")[0]
			accountsList = append(accountsList, name)
		}
	}

	return accountsList
}

// RetrievalAttempts returns the configured retrieval attempts or the default
// value if not specified.
func (c Config) RetrievalAttempts() int {
	if c.FetcherOAuth2TokenSettings.RetrievalAttempts <= 0 {
		return defaultTokenRetrievalAttempts
	}
	return c.FetcherOAuth2TokenSettings.RetrievalAttempts
}
