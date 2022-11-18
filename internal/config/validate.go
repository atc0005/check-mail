// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"strings"
)

func validateTLSVersion(c Config) error {
	switch strings.ToLower(c.minTLSVersion) {
	case minTLSVersion10:
		return nil
	case minTLSVersion11:
		return nil
	case minTLSVersion12:
		return nil
	case minTLSVersion13:
		return nil
	default:
		return fmt.Errorf("invalid TLS version keyword: %s", c.minTLSVersion)
	}
}

func validateNetworkType(c Config) error {
	switch strings.ToLower(c.NetworkType) {
	case netTypeTCPAuto:
		return nil
	case netTypeTCP4:
		return nil
	case netTypeTCP6:
		return nil
	default:
		return fmt.Errorf("invalid network type keyword: %s", c.NetworkType)
	}
}

func validateLoggingLevels(c Config) error {
	requestedLoggingLevel := strings.ToLower(c.LoggingLevel)
	if _, ok := loggingLevels[requestedLoggingLevel]; !ok {
		return fmt.Errorf("invalid logging level %s", c.LoggingLevel)
	}

	return nil
}

func validateAccounts(c Config) error {
	for _, account := range c.Accounts {
		if account.Folders == nil {
			return fmt.Errorf(
				"one or more folders not provided for account %s",
				account.Name,
			)
		}

		if account.Port < 0 {
			return fmt.Errorf(
				"invalid TCP port number %d provided for account %s",
				account.Port,
				account.Name,
			)
		}

		if account.Username == "" {
			return fmt.Errorf("username not provided for account %s",
				account.Name,
			)
		}

		if account.Password == "" {
			return fmt.Errorf("password not provided for account %s",
				account.Name,
			)
		}

		if account.Server == "" {
			return fmt.Errorf("server FQDN not provided for account %s",
				account.Name,
			)
		}
	}

	return nil
}

// validate verifies all Config struct fields have been provided acceptable
// values.
func (c Config) validate(appType AppType) error {

	switch {
	case appType.PluginIMAPMailboxBasicAuth:

		if err := validateAccounts(c); err != nil {
			return err
		}

		if err := validateTLSVersion(c); err != nil {
			return err
		}

		if err := validateNetworkType(c); err != nil {
			return err
		}

		if err := validateLoggingLevels(c); err != nil {
			return err
		}

	case appType.ReporterIMAPMailboxBasicAuth:

		// NOTE: It's fine to *not* specify a config file. The expected behavior
		// is that specifying a config file will be a rare thing; users will more
		// often than not rely on config file auto-detection behavior.
		//
		// That said, if a user does not specify a config file, we need to require
		// that one was found and loaded.
		//
		// if useConfigFile {
		// 	if c.ConfigFile == "" {
		// 		return fmt.Errorf("config file required, but not specified")
		// 	}
		// }

		// set with a default value if not specified by the user, so should not
		// ever be empty
		if c.ReportFileOutputDir == "" {
			return fmt.Errorf("missing report file output directory")
		}

		// set with a default value if not specified by the user, so should not
		// ever be empty
		if c.LogFileOutputDir == "" {
			return fmt.Errorf("missing log file output directory")
		}

		if err := validateAccounts(c); err != nil {
			return err
		}

		if err := validateTLSVersion(c); err != nil {
			return err
		}

		if err := validateNetworkType(c); err != nil {
			return err
		}

		if err := validateLoggingLevels(c); err != nil {
			return err
		}

	default:
		return fmt.Errorf(
			"unable to validate configuration: %w",
			ErrAppTypeNotSpecified,
		)
	}

	// Optimist
	return nil
}
