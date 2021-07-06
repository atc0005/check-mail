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
