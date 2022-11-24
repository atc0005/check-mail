// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

const (

	// LogLevelDisabled maps to zerolog.Disabled logging level
	LogLevelDisabled string = "disabled"

	// LogLevelPanic maps to zerolog.PanicLevel logging level
	LogLevelPanic string = "panic"

	// LogLevelFatal maps to zerolog.FatalLevel logging level
	LogLevelFatal string = "fatal"

	// LogLevelError maps to zerolog.ErrorLevel logging level
	LogLevelError string = "error"

	// LogLevelWarn maps to zerolog.WarnLevel logging level
	LogLevelWarn string = "warn"

	// LogLevelInfo maps to zerolog.InfoLevel logging level
	LogLevelInfo string = "info"

	// LogLevelDebug maps to zerolog.DebugLevel logging level
	LogLevelDebug string = "debug"

	// LogLevelTrace maps to zerolog.TraceLevel logging level
	LogLevelTrace string = "trace"
)

// LoggingLevels is a map of string to zerolog.Level created in an effort to
// keep from repeating ourselves
var loggingLevels = make(map[string]zerolog.Level)

func init() {

	// https://stackoverflow.com/a/59426901
	// syntax error: non-declaration statement outside function body
	//
	// Workaround: Use init() to setup this map for later reference
	loggingLevels[LogLevelDisabled] = zerolog.Disabled
	loggingLevels[LogLevelPanic] = zerolog.PanicLevel
	loggingLevels[LogLevelFatal] = zerolog.FatalLevel
	loggingLevels[LogLevelError] = zerolog.ErrorLevel
	loggingLevels[LogLevelWarn] = zerolog.WarnLevel
	loggingLevels[LogLevelInfo] = zerolog.InfoLevel
	loggingLevels[LogLevelDebug] = zerolog.DebugLevel
	loggingLevels[LogLevelTrace] = zerolog.TraceLevel
}

// setLoggingLevel applies the requested logging level to filter out messages
// with a lower level than the one configured.
func setLoggingLevel(logLevel string) error {

	switch logLevel {
	case LogLevelDisabled:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelDisabled])
	case LogLevelPanic:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelPanic])
	case LogLevelFatal:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelFatal])
	case LogLevelError:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelError])
	case LogLevelWarn:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelWarn])
	case LogLevelInfo:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelInfo])
	case LogLevelDebug:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelDebug])
	case LogLevelTrace:
		zerolog.SetGlobalLevel(loggingLevels[LogLevelTrace])
	default:
		return fmt.Errorf("invalid option provided: %v", logLevel)
	}

	// signal that a case was triggered as expected
	return nil

}

// setupLogging is responsible for configuring logging settings for this
// application
func (c *Config) setupLogging(appType AppType) error {

	switch {

	// we want to log to a file only for list-emails
	case appType.ReporterIMAPMailbox:

		logFilename := fmt.Sprintf(
			logFilenameTemplate,
			time.Now().Format(logFilenameDateLayout),
		)

		if err := os.MkdirAll(c.LogFileOutputDir, defaultDirectoryPerms); err != nil {
			return fmt.Errorf("failed to create log output dir: %w", err)
		}

		logFilePath := filepath.Join(c.LogFileOutputDir, logFilename)

		f, fileOpenErr := os.OpenFile(
			filepath.Clean(logFilePath),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			defaultFilePerms,
		)

		if fileOpenErr != nil {
			return fmt.Errorf("failed to open report file: %w", fileOpenErr)
		}

		// TODO: What is a better way to ensure this is closed properly?
		// Currently this is closed from main() as a deferred call.
		c.LogFileHandle = f

		// *Nearly* everything for this app type is sent to the log file for
		// later review. We use the "console writer" in an effort to make the
		// log file easier to visually review.
		logOutput := zerolog.ConsoleWriter{Out: f, NoColor: true}

		c.Log = zerolog.New(logOutput).With().Timestamp().Caller().
			Str("version", Version()).
			Str("network_type", c.NetworkType).
			Str("min_tls_version", c.MinTLSVersionKeyword()).
			Logger()

	case appType.InspectorIMAPCaps:

		// Slimline logger to emit messages in a format more appropriate to
		// CLI "inspector" tool.
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		c.Log = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()

	case appType.PluginIMAPMailboxBasicAuth:

		// Whatever output meant for consumption is emitted to stdout and
		// whatever is meant for troubleshooting is sent to stderr. To help
		// keep these two goals separate (and because Nagios doesn't really do
		// anything special with JSON output from plugins), we use stdlib fmt
		// package output functions for Nagios via stdout and logging package
		// for troubleshooting via stderr.
		//
		// If we're not setting up the configuration for the Nagios plugin, we
		// will attempt to use another output target.
		logOutput := os.Stderr

		consoleWriter := zerolog.ConsoleWriter{Out: logOutput, NoColor: true}
		c.Log = zerolog.New(consoleWriter).With().Timestamp().Caller().
			Str("version", Version()).
			Str("network_type", c.NetworkType).
			Str("min_tls_version", c.MinTLSVersionKeyword()).
			Logger()

	case appType.PluginIMAPMailboxOAuth2:

		// Whatever output meant for consumption is emitted to stdout and
		// whatever is meant for troubleshooting is sent to stderr. To help
		// keep these two goals separate (and because Nagios doesn't really do
		// anything special with JSON output from plugins), we use stdlib fmt
		// package output functions for Nagios via stdout and logging package
		// for troubleshooting via stderr.
		//
		// If we're not setting up the configuration for the Nagios plugin, we
		// will attempt to use another output target.
		logOutput := os.Stderr

		consoleWriter := zerolog.ConsoleWriter{Out: logOutput, NoColor: true}
		c.Log = zerolog.New(consoleWriter).With().Timestamp().Caller().
			Str("version", Version()).
			Str("network_type", c.NetworkType).
			Str("min_tls_version", c.MinTLSVersionKeyword()).
			Logger()

	}

	if err := setLoggingLevel(c.LoggingLevel); err != nil {
		return err
	}

	return nil

}
