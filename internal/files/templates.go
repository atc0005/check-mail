// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package files

// ReportFilenameTemplate is used as a filename template when generating new
// email check report files. This file is appended to for each account
// processed. In order to help prevent repeatedly writing to the same report
// the current date/time when first writing to the file is used as part of the
// filename itself.
const ReportFilenameTemplate string = "%s-emails-report-%s.txt"

// ReportFilenameDateLayout is the time formatting layout for generated
// reports.
const ReportFilenameDateLayout string = "2006-01-02-15_04"

// Mon Jan 2 15:04:05 -0700 MST 2006

// reportFileTemplateText is used when generating the output report file. This
// file is provided to the team as part of the daily email checks (e.g.,
// software updates, vulnerability announcements, etc.).
const reportFileTemplateText string = `
h4. {{ .AccountName }}

h5. Overview

| Report generated | {{ .ReportTime.Format "2006-01-02 15:04:05" }} |
| Summary | {{ .MessagesFoundSummary }} |
| Placeholder character | {{ .UnicodeCharSubstitute }} (substituted for Emoji incompatible with MySQL utf8mb3 character set) |

h5. Emails found

|_.Folder|_.Subject|_.Date|
{{ range .MailboxCheckResults -}}
{{- $mailboxName := .MailboxName -}}
{{- range .Messages -}}
{{- if .ModifiedSubject -}}
| {{ $mailboxName }} | {{ .ModifiedSubject }} | {{ .EnvelopeDateFormatted }} |
{{- else -}}
| {{ $mailboxName }} | {{ .OriginalSubject }} | {{ .EnvelopeDateFormatted }} |
{{- end }}
{{ end -}}
{{ end }}

h5. Reported emails

|_.Folder|_.Subject|_.Date|_.Reported on|
| PLACEHOLDER | PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |

h5. Emails moved to applicable folders

|_.Folder|_.Subject|_.Date|
| PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |

h5. Deleted emails

|_.Folder|_.Subject|_.Date|
| PLACEHOLDER | PLACEHOLDER | PLACEHOLDER |
`
