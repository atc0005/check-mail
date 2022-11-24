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

// validateTLSVersion asserts that the specified TLS version keyword is valid.
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

// validateNetworkType asserts that the requested network type keyword is
// valid.
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

// validateLoggingLevels asserts that the requested logging level is valid.
func validateLoggingLevels(c Config) error {
	requestedLoggingLevel := strings.ToLower(c.LoggingLevel)
	if _, ok := loggingLevels[requestedLoggingLevel]; !ok {
		return fmt.Errorf("invalid logging level %s", c.LoggingLevel)
	}

	return nil
}

// validateAccountBasicAuthFields is responsible for validating MailAccount
// fields specific to the Basic Authentication type. The caller is responsible
// for calling this function for the appropriate application type.
func validateAccountBasicAuthFields(account MailAccount, appType AppType) error {
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

	return nil
}

// validateAccountOAuth2ClientCredsAuthFields is responsible for validating
// MailAccount fields specific to the OAuth2 Client Credentials Flow
// authentication type. The caller is responsible for calling this function
// for the appropriate application type.
func validateAccountOAuth2ClientCredsAuthFields(account MailAccount, appType AppType) error {
	if account.OAuth2Settings.ClientID == "" {
		return fmt.Errorf("client ID not provided for account %s",
			account.Name,
		)
	}

	if account.OAuth2Settings.ClientSecret == "" {
		return fmt.Errorf("client secret not provided for account %s",
			account.Name,
		)
	}

	// Scopes is non-optional. If we want to support just *one* IMAP provider
	// (e.g., O365) we can fallback to using a default scope, but if the goal
	// is (it is) to support multiple providers we need to require at least
	// one scope value.
	if len(account.OAuth2Settings.Scopes) == 0 {
		return fmt.Errorf("scopes not provided for account %s",
			account.Name,
		)
	}

	// Unlikely to have empty slice strings, but worth ruling out?
	for _, scope := range account.OAuth2Settings.Scopes {
		if strings.TrimSpace(scope) == "" {
			return fmt.Errorf("empty scope provided for account %s",
				account.Name,
			)
		}
	}

	if account.OAuth2Settings.SharedMailbox == "" {
		return fmt.Errorf("shared mailbox name not provided for account %s",
			account.Name,
		)
	}

	if account.OAuth2Settings.TokenURL == "" {
		return fmt.Errorf("token URL not provided for account %s",
			account.Name,
		)
	}

	return nil
}

// validateAccounts is responsible for validating MailAccount fields.
func validateAccounts(c Config, appType AppType) error {
	for _, account := range c.Accounts {

		// All app types use this field.
		if account.Port < 0 {
			return fmt.Errorf(
				"invalid TCP port number %d provided for account %s",
				account.Port,
				account.Name,
			)
		}

		// All app types use this field.
		if account.Server == "" {
			return fmt.Errorf("server FQDN not provided for account %s",
				account.Name,
			)
		}

		switch {
		case appType.ReporterIMAPMailbox:

			switch account.AuthType {
			case AuthTypeBasic:
				if err := validateAccountBasicAuthFields(account, appType); err != nil {
					return err
				}

			case AuthTypeOAuth2ClientCreds:
				if err := validateAccountOAuth2ClientCredsAuthFields(account, appType); err != nil {
					return err
				}

			default:
				return fmt.Errorf(
					"unexpected authentication type %q: %w",
					account.AuthType,
					ErrInvalidAuthType,
				)
			}

			if account.Folders == nil {
				return fmt.Errorf(
					"one or more folders not provided for account %s",
					account.Name,
				)
			}

		case appType.InspectorIMAPCaps:

			// This app type only uses the server/port values.

		case appType.PluginIMAPMailboxBasicAuth:
			if err := validateAccountBasicAuthFields(account, appType); err != nil {
				return err
			}

			if account.Folders == nil {
				return fmt.Errorf(
					"one or more folders not provided for account %s",
					account.Name,
				)
			}

		case appType.PluginIMAPMailboxOAuth2:

			if err := validateAccountOAuth2ClientCredsAuthFields(account, appType); err != nil {
				return err
			}

			if account.Folders == nil {
				return fmt.Errorf(
					"one or more folders not provided for account %s",
					account.Name,
				)
			}
		}

	}

	return nil
}

// validate verifies all Config struct fields have been provided acceptable
// values.
func (c Config) validate(appType AppType) error {

	switch {
	case appType.InspectorIMAPCaps:

		// We're using the Accounts collection in order to obtain access to
		// the server and port fields.
		if err := validateAccounts(c, appType); err != nil {
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

	case appType.PluginIMAPMailboxBasicAuth:

		if err := validateAccounts(c, appType); err != nil {
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

	case appType.PluginIMAPMailboxOAuth2:

		if err := validateAccounts(c, appType); err != nil {
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

	case appType.ReporterIMAPMailbox:

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

		if err := validateAccounts(c, appType); err != nil {
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
