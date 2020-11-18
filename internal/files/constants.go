// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package files

import "os"

// default permissions granting owner full access, deny access to all others
const (
	defaultDirectoryPerms os.FileMode = 0700
	defaultFilePerms      os.FileMode = 0600
)

// ReportFileOpenFlags is configured to create report files in append-only
// write mode. This works well when writing all report data to a single file
// and works fairly well when writing to separate files, one per account name.
// One potential side-effect is that multiple executions of the application
// can result in appended reports instead of separate reports due to how we
// are currently using date/time based naming patterns. If this becomes a
// problem, we should update this flag set to overwrite any existing contents
// so that only the latest execution of the application (within the same
// minute) is reflected within the report file.
const ReportFileOpenFlags int = os.O_APPEND | os.O_CREATE | os.O_WRONLY
