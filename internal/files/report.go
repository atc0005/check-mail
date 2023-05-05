// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/atc0005/check-mail/internal/mbxs"
	"github.com/rs/zerolog"
)

// ReportData represents the values that are used when generating reports via
// templates.
type ReportData struct {

	// AccountName is the name of the email account that was processed. This
	// matches the section name used in the configuration file.
	AccountName string

	// MailboxCheckResults provides the final state of our email accounts
	// evaluation that we will use for generating our reports.
	MailboxCheckResults mbxs.MailboxCheckResults

	// MessagesFoundSummary is a one-line summary of the mail items found in
	// checked mailboxes.
	MessagesFoundSummary string

	// ReportTime is when the email summary report was generated.
	ReportTime time.Time

	// String used in place of incompatible Unicode characters for the utf8mb3
	// character set.
	UnicodeCharSubstitute string

	// What else to include here?
}

// GenerateReport is a wrapper around the steps needed for generating or
// updating a email summary report. This function receives a ReportData type
// that acts as a container for all required information used by the report
// generation process and the path to a directory that will hold generated
// reports.
func GenerateReport(reportData ReportData, reportDirectory string, logger zerolog.Logger) error {

	// very unlikely that we'll need to worry with concurrency when working
	// with this file, but playing it safe?
	var mutex = &sync.Mutex{}

	reportFilename := fmt.Sprintf(
		ReportFilenameTemplate,
		reportData.AccountName,
		reportData.ReportTime.Format(ReportFilenameDateLayout),
	)

	reportFilePath := filepath.Join(reportDirectory, reportFilename)

	logger = logger.With().
		Str("report_file", reportFilePath).
		Str("report_directory", reportDirectory).
		Logger()

	reportFileTemplate, tmplParseErr := template.
		New("reportFileTemplate").Parse(reportFileTemplateText)
	if tmplParseErr != nil {
		return fmt.Errorf(
			"failed to parse report file template: %w",
			tmplParseErr,
		)
	}

	logger.Debug().Msg("Calling os.MkdirAll to ensure report directory exists")

	if err := os.MkdirAll(reportDirectory, defaultDirectoryPerms); err != nil {
		return fmt.Errorf("failed to create report output dir: %w", err)
	}

	logger.Debug().Msg("No errors calling os.MkdirAll")

	// Append to the file, not overwrite it as we will be updating the file
	// for each account evaluated.
	f, fileOpenErr := os.OpenFile(
		filepath.Clean(reportFilePath),
		ReportFileOpenFlags,
		defaultFilePerms,
	)
	if fileOpenErr != nil {
		return fmt.Errorf("failed to open report file: %w", fileOpenErr)
	}

	logger.Debug().Msg("Successfully opened report file for updates")

	// #nosec G307
	// Believed to be a false-positive from recent gosec release
	// https://github.com/securego/gosec/issues/714
	defer func(filename string) {
		if err := f.Close(); err != nil {
			// Ignore "file already closed" errors
			if !errors.Is(err, os.ErrClosed) {
				logger.Error().
					Err(err).
					Str("filename", filename).
					Msg("failed to close file")
			}
		}
	}(reportFilePath)

	logger.Debug().Msg("Locking mutex")
	mutex.Lock()
	defer func() {
		logger.Debug().Msg("Unlocking mutex")
		mutex.Unlock()
	}()

	logger.Debug().Msg("Executing template to update report file")
	if tmplErr := reportFileTemplate.Execute(f, reportData); tmplErr != nil {

		// if there were template execution errors, go ahead and try to close
		// the file before returning the template write error
		if fileCloseErr := f.Close(); fileCloseErr != nil {

			// log this error, return Write error as it takes precedence
			logger.Error().Err(fileCloseErr).Msg("failed to close report file")
		}

		return fmt.Errorf(
			"error writing to file %q: %w",
			reportFilePath,
			tmplErr,
		)
	}
	logger.Debug().Msg("Successfully executed template to update report file")

	logger.Debug().Msg("Syncing file modifications")
	if syncErr := f.Sync(); syncErr != nil {
		return fmt.Errorf(
			"failed to explicitly sync file %q after writing: %w",
			reportFilePath,
			syncErr,
		)
	}
	logger.Debug().Msg("Successfully synced modifications to report file")

	logger.Debug().Msg("Closing report file")
	if closeErr := f.Close(); closeErr != nil {
		return fmt.Errorf(
			"error closing file %q: %w",
			reportFilePath,
			closeErr,
		)
	}
	logger.Debug().Msg("Successfully closed report file")

	return nil
}
