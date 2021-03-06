// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// readConfigFile receives a list of potential configuration files, opens the
// first available configuration file from the list and returns the contents
// for parsing. The config file that was successfully read is recorded for
// later reference. An error is returned if no listed configuration file could
// be read.
func (c *Config) readConfigFile(configFile ...string) ([]byte, error) {

	var data []byte
	var finalErr error

	for _, file := range configFile {
		c.Log.Debug().Str("config_file", file).Msg("Attempting to open config file")

		fh, err := os.Open(filepath.Clean(file))
		if err != nil {
			// Encountering an error here may not be fatal (file may not be
			// present). Record the error and try the next one.
			finalErr = err
			continue
		}
		c.Log.Debug().Str("config_file", file).Msg("Config file opened")
		defer func() {
			if err := fh.Close(); err != nil {
				// Ignore "file already closed" errors
				if !errors.Is(err, os.ErrClosed) {
					c.Log.Error().Msgf(
						"failed to close file %q: %s",
						file,
						err.Error(),
					)
				}
			}
		}()

		c.Log.Debug().Msg("Attempting to read config file")

		var readErr error
		data, readErr = ioutil.ReadAll(fh)
		if readErr != nil {
			finalErr = fmt.Errorf("failed to read config file: %w", readErr)
			continue
		}

		// We just successfully read a configuration file for later parsing.
		// Clear any errors that may have been recorded earlier (e.g., from
		// attempting to open non-existent files) so that we don't
		// unintentionally consider those transient errors to be permanent
		// errors.
		finalErr = nil

		// Note which configuration file was successfully read
		c.ConfigFileUsed = file

		// stop processing further config files on the first successful read
		break
	}

	// After attempting to open/process all provided config files, see if we
	// ended with an error. If so, consider the config file read attempt (of
	// all provided files) an error.
	if finalErr != nil {
		return nil, fmt.Errorf("failed to read config file: %w", finalErr)
	}

	// If we made it this far, we were able to successfully read one of the
	// provided config files, so we consider the attempt a success.
	return data, nil
}

// loadConfigFile is a helper function to handle opening a config file and
// importing the settings for use. Multiple config file paths are accepted and
// tried in order. The first config file successfully loaded is recorded for
// later reference. A successful load aborts attempts to process any remaining
// config files in the list. If no config file can be successfully loaded an
// error is returned.
func (c *Config) loadConfigFile(configFile ...string) error {

	data, readErr := c.readConfigFile(configFile...)
	if readErr != nil {
		c.ConfigFileUsed = ""
		return fmt.Errorf(
			"failed to load config file %v: %w",
			configFile,
			readErr,
		)
	}

	loadErr := c.loadINIConfig(data)
	if loadErr != nil {
		c.ConfigFileUsed = ""
		return fmt.Errorf(
			"failed to load configuration from INI file %q: %w",
			c.ConfigFileUsed,
			loadErr,
		)
	}

	// Explicitly note that a config file *was* loaded. The config file read
	// is recorded as c.ConfigFileUsed.
	c.ConfigFileLoaded = true

	return nil

}

// localConfigFile returns the potential path to a config file found locally
// alongside the executable, or an error if one is encountered constructing
// the path.
func (c Config) localConfigFile(filename string) (string, error) {

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf(
			"unable to get running executable path to load local config file: %w",
			err,
		)
	}
	exeDirPath, _ := filepath.Split(exePath)
	localCfgFile := filepath.Join(exeDirPath, filename)

	c.Log.Debug().Msgf("local config file path: %q", localCfgFile)

	return localCfgFile, nil
}

// userConfigFile returns the potential path to a config file found in the
// user's configuration path or an error if one is encountered constructing the
// path.
func (c Config) userConfigFile(appName string, filename string) (string, error) {
	// Ubuntu environment:
	// os.UserHomeDir: /home/username
	// os.UserConfigDir: /home/username/.config
	//
	// Windows environment:
	// os.UserHomeDir: C:\Users\username
	// os.UserConfigDir: C:\Users\username\AppData\Roaming

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user config dir: %w", err)
	}

	userConfigAppDir := filepath.Join(userConfigDir, appName)
	userConfigFileFullPath := filepath.Join(userConfigAppDir, filename)
	c.Log.Debug().Msgf("user config file path: %q", userConfigFileFullPath)

	return userConfigFileFullPath, nil
}

// loadINIConfig loads or imports the contents of an INI configuration file
// and populates the list of accounts used by applications in this project.
// Any error encountered during the import process is returned.
func (c *Config) loadINIConfig(file []byte) error {

	iniFile, loadErr := ini.Load(file)
	if loadErr != nil {
		return fmt.Errorf("failed to load INI file: %v", loadErr)
	}

	defaultSection, lookupErr := iniFile.GetSection("DEFAULT")
	if lookupErr != nil {
		return fmt.Errorf(
			"failed to retrieve defaults section: %v",
			lookupErr,
		)
	}

	c.Log.Debug().Msgf(
		"Keys for section %s: %+v\n",
		defaultSection.Name(),
		defaultSection.Keys(),
	)

	serverNameKey, lookupErr := defaultSection.GetKey(iniDefaultServerNameKeyName)
	if lookupErr != nil {
		return fmt.Errorf(
			"failed to retrieve value from key %s in section %s",
			iniDefaultServerNameKeyName,
			defaultSection.Name(),
		)
	}
	serverName := serverNameKey.Value()

	serverPortKey, lookupErr := defaultSection.GetKey(iniDefaultServerPortKeyName)
	if lookupErr != nil {
		return fmt.Errorf(
			"failed to retrieve value from key %s: %v",
			iniDefaultServerPortKeyName,
			lookupErr,
		)
	}

	// convert string to int
	serverPort, strConvErr := strconv.Atoi(serverPortKey.Value())
	if strConvErr != nil {
		return fmt.Errorf(
			"failed to convert string %q to int: %v",
			serverPortKey.Value(),
			strConvErr,
		)
	}

	// at this point we already have serverName, serverPort values. We
	// need to retrieve the account-specific values provided by each
	// unique section in the INI file.
	for _, section := range iniFile.Sections() {
		if strings.ToLower(section.Name()) == "default" {
			c.Log.Debug().Msg("found/skipping lowercase default section")
			continue
		}

		// Not used directly for much at the moment, but it will be needed to
		// generate the replacement TOML config file later.
		accountName := section.Name()

		usernameKey, lookupErr := section.GetKey(iniUsernameKeyName)
		if lookupErr != nil {
			return fmt.Errorf(
				"failed to retrieve value from key %s: %v",
				usernameKey,
				lookupErr,
			)
		}
		username := usernameKey.Value()

		passwordKey, lookupErr := section.GetKey(iniPasswordKeyName)
		if lookupErr != nil {
			return fmt.Errorf(
				"failed to retrieve value from key %s: %v",
				iniPasswordKeyName,
				lookupErr,
			)
		}
		password := passwordKey.Value()

		foldersKey, lookupErr := section.GetKey(iniFoldersKeyName)
		if lookupErr != nil {
			return fmt.Errorf(
				"failed to retrieve value from key %s: %v",
				iniFoldersKeyName,
				lookupErr,
			)
		}

		// split and trim folders list provided as single string in INI file.
		folders := strings.Split(foldersKey.Value(), ",")
		for i, folder := range folders {
			folders[i] = strings.Trim(folder, `" `)
		}

		account := MailAccount{
			Server:   serverName,
			Port:     serverPort,
			Name:     accountName,
			Username: username,
			Password: password,
			Folders:  folders,
		}
		c.Accounts = append(c.Accounts, account)

	}

	return nil
}
