// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package logging provides a logging configuration wrapper for the
// applications within this repo
package logging

import (
	"fmt"

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
var LoggingLevels = make(map[string]zerolog.Level)

func init() {

	// https://stackoverflow.com/a/59426901
	// syntax error: non-declaration statement outside function body
	//
	// Workaround: Use init() to setup this map for later reference
	LoggingLevels[LogLevelDisabled] = zerolog.Disabled
	LoggingLevels[LogLevelPanic] = zerolog.PanicLevel
	LoggingLevels[LogLevelFatal] = zerolog.FatalLevel
	LoggingLevels[LogLevelError] = zerolog.ErrorLevel
	LoggingLevels[LogLevelWarn] = zerolog.WarnLevel
	LoggingLevels[LogLevelInfo] = zerolog.InfoLevel
	LoggingLevels[LogLevelDebug] = zerolog.DebugLevel
	LoggingLevels[LogLevelTrace] = zerolog.TraceLevel
}

// SetLoggingLevel applies the requested logging level to filter out messages
// with a lower level than the one configured.
func SetLoggingLevel(logLevel string) error {

	switch logLevel {
	case LogLevelDisabled:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelDisabled])
	case LogLevelPanic:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelPanic])
	case LogLevelFatal:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelFatal])
	case LogLevelError:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelError])
	case LogLevelWarn:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelWarn])
	case LogLevelInfo:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelInfo])
	case LogLevelDebug:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelDebug])
	case LogLevelTrace:
		zerolog.SetGlobalLevel(LoggingLevels[LogLevelTrace])
	default:
		return fmt.Errorf("invalid option provided: %v", logLevel)
	}

	// signal that a case was triggered as expected
	return nil

}
