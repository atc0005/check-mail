<!-- omit in toc -->
# Check Mail

Various tools used to monitor mail services

[![Latest Release](https://img.shields.io/github/release/atc0005/check-mail.svg?style=flat-square)](https://github.com/atc0005/check-mail/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-mail.svg)](https://pkg.go.dev/github.com/atc0005/check-mail)
[![Validate Codebase](https://github.com/atc0005/check-mail/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/check-mail/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/check-mail/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/check-mail/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/check-mail/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/check-mail/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/check-mail/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/check-mail/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of contents

- [Project home](#project-home)
- [Overview](#overview)
- [Features](#features)
  - [`check_imap_mailbox`](#check_imap_mailbox)
  - [`list-emails`](#list-emails)
- [Requirements](#requirements)
- [Installation](#installation)
  - [From source](#from-source)
  - [Using release binaries](#using-release-binaries)
- [Configuration Options](#configuration-options)
  - [`check_imap_mailbox`](#check_imap_mailbox-1)
    - [Command-line arguments](#command-line-arguments)
  - [`list-emails`](#list-emails-1)
    - [Command-line arguments](#command-line-arguments-1)
    - [Configuration file](#configuration-file)
      - [Settings](#settings)
      - [Usage](#usage)
- [Examples](#examples)
  - [`check_imap_mailbox`](#check_imap_mailbox-2)
    - [As a Nagios plugin](#as-a-nagios-plugin)
    - [Login failure](#login-failure)
    - [Help Output](#help-output)
  - [`list-emails`](#list-emails-2)
    - [No options](#no-options)
    - [Alternate locations for config file, log and report directories](#alternate-locations-for-config-file-log-and-report-directories)
    - [Help output](#help-output-1)
- [License](#license)
- [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/check-mail) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

This repo contains various tools used to monitor mail services.

| Tool Name            | Overall Status | Description                                                |
| -------------------- | -------------- | ---------------------------------------------------------- |
| `check_imap_mailbox` | Stable         | Nagios plugin used to monitor mailboxes for items          |
| `list-emails`        | Beta           | Small CLI app used to generate listing of mailbox contents |

## Features

### `check_imap_mailbox`

- Monitor specified mailboxes for an IMAP account
  - non-`OK` state returned for any items in specified mailboxes or errors
    encountered
  - `OK` state returned if all specified mailboxes are empty
- Leveled logging
  - JSON-format output
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`
- TLS IMAP4 connectivity
  - port defaults to 993/tcp
  - network type defaults to either of IPv4 and IPv6, but optionally limited
    to IPv4-only or IPv6-only
- Optional branding "signature"
  - used to indicate what Nagios plugin (and what version) is responsible for
    the service check result

### `list-emails`

- Check one or many mailboxes
- Leveled logging
  - JSON-format output
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`
  - bulk of logging directed to per-invocation log file
- TLS IMAP4 connectivity
  - port defaults to 993/tcp
  - network type defaults to either of IPv4 and IPv6, but optionally limited
    to IPv4-only or IPv6-only
- Minimal output to console unless requested
  - via `debug` logging level
- Textile (Redmine compatible) formatted report generated per specified email
  account
  - overall summary
  - template copy/paste/modify report for posting to Redmine issues (aka,
    "tickets")
  - automatic replacement of Unicode characters outside of the MySQL `utf8mb3`
    character set with a placeholder character
    - the intent is to help prevent MySQL errors when posting summary reports
      - e.g., `ERROR 1366 (22007): Incorrect string value`

## Requirements

- Go
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

Tested using:

- Go 1.13+
- Windows 10 Version 1903
  - native
  - WSL
- Ubuntu Linux 16.04+

## Installation

### From source

1. [Download](https://golang.org/dl/) Go
1. [Install](https://golang.org/doc/install) Go
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/check-mail`
   1. `cd check-mail`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     1. `sudo yum install make gcc`
1. Build
   - for current operating system
     - `go build -mod=vendor ./cmd/check_imap_mailbox/`
       - *forces build to use bundled dependencies in top-level `vendor`
         folder*
     - `go build -mod=vendor ./cmd/list-emails/`
       - *forces build to use bundled dependencies in top-level `vendor`
         folder*
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Locate generated binaries
   - if using `Makefile`
     - look in `/tmp/check-mail/release_assets/check_imap_mailbox/`
     - look in `/tmp/check-mail/release_assets/list-emails/`
   - if using `go build`
     - look in `/tmp/check-mail/`
1. Copy the applicable binaries to whatever systems needs to run them
1. Deploy
   - Place `list-emails` in a location of your choice
   - Place `check_imap_mailbox` in the same location where your distro's
     package manage has place other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_imap_mailbox` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_imap_mailbox` on RedHat-based
       systems
1. Copy the template [configuration file](#configuration-file), modify
   accordingly and place in a [supported location](#configuration-file)

### Using release binaries

1. Download the [latest
   release](https://github.com/atc0005/check-mail/releases/latest) binaries
1. Deploy
   - Place `list-emails` in a location of your choice
   - Place `check_imap_mailbox` in the same location where your distro's
     package manager places other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_imap_mailbox` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_imap_mailbox` on RedHat-based
       systems
1. Copy the template [configuration file](#configuration-file), modify
   accordingly and place in a [supported location](#configuration-file)

## Configuration Options

### `check_imap_mailbox`

#### Command-line arguments

- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined.

| Option          | Required | Default        | Repeat | Possible                                                                | Description                                                                                                                                                                                 |
| --------------- | -------- | -------------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       |                | No     | `-h`, `--help`                                                          | Generate listing of all valid command-line options and applicable (short) guidance for using them.                                                                                          |
| `folders`       | Yes      | *empty string* | No     | *comma-separated list of folders*                                       | Folders or IMAP "mailboxes" to check for mail. This value is provided as a comma-separated list.                                                                                            |
| `username`      | Yes      | *empty string* | No     | *valid username, often in email address format*                         | The account used to login to the remote mail server. This is often in the form of an email address.                                                                                         |
| `password`      | Yes      | *empty string* | No     | *valid password*                                                        | The remote mail server account password.                                                                                                                                                    |
| `server`        | Yes      | *empty string* | No     | *valid FQDN or IP Address*                                              | The fully-qualified domain name of the remote mail server.                                                                                                                                  |
| `port`          | No       | `993`          | No     | *valid IMAP TCP port*                                                   | TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections.                                                                  |
| `net-type`      | No       | `auto`         | No     | `auto`, `tcp4`, `tcp6`                                                  | Limits network connections to remote mail servers to one of the specified types.                                                                                                            |
| `logging-level` | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                                                                                             |
| `branding`      | No       | `false`        | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. Because this output may not mix well with branding information emitted by other tools, this output is disabled by default. |
| `version`       | No       | `false`        | No     | `true`, `false`                                                         | Whether to display application version and then immediately exit application                                                                                                                |

### `list-emails`

#### Command-line arguments

- The bulk of the settings for this application are provided via the
  `accounts.ini` configuration file.
- It is not currently possible to specify all required settings by
  command-line

| Option            | Required | Default        | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                              |
| ----------------- | -------- | -------------- | ------ | ----------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`       | No       |                | No     | `-h`, `--help`                                                          | Generate listing of all valid command-line options and applicable (short) guidance for using them.                                                                                                                                                                                                                                       |
| `config-file`     | No       | `accounts.ini` | No     | *valid path to INI configuration file for this application*             | Full path to the INI-formatted configuration file used by this application. See contrib/list-emails/accounts.example.ini for a starter template. Rename to accounts.ini, update with applicable information and place in a directory of your choice. If this file is found in your current working directory you need not use this flag. |
| `log-file-dir`    | No       | `log`          | No     | *valid, writable path to a directory*                                   | Full path to the directory where log files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist.                                                                 |
| `report-file-dir` | No       | `output`       | No     | *valid, writable path to a directory*                                   | Full path to the directory where email summary report files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist.                                                |
| `net-type`        | No       | `auto`         | No     | `auto`, `tcp4`, `tcp6`                                                  | Limits network connections to remote mail servers to one of the specified types.                                                                                                                                                                                                                                                         |
| `logging-level`   | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                                                                                                                                                                                                                                          |
| `version`         | No       | `false`        | No     | `true`, `false`                                                         | Whether to display application version and then immediately exit application                                                                                                                                                                                                                                                             |

#### Configuration file

##### Settings

**NOTE**: The `email1` and `email2` values below are for illustration. You are
free to choose section names, though it is recommended to base them off of the
username (sans `@` symbol and domain part) for each email account. While
`email1` and `email2` are both listed, only one is required (to reflect one
account) though many such entries (one per account) are supported.

| Config file Setting Name | Section Name | Notes                                               |
| ------------------------ | ------------ | --------------------------------------------------- |
| `server_name`            | `DEFAULT`    | FQDN of IMAP server (e.g., `outlook.office365.com`) |
| `server_port`            | `DEFAULT`    | Usually 993                                         |
| `username`               | `email1`     | Often in the form of an email address               |
| `password`               | `email1`     | Account password                                    |
| `folders`                | `email1`     | Double quoted, comma separated                      |
| `username`               | `email2`     | Often in the form of an email address               |
| `password`               | `email2`     | Account password                                    |
| `folders`                | `email2`     | Double quoted, comma separated                      |

##### Usage

The [accounts.example.ini](contrib/list-emails/accounts.example.ini) INI
config file is intended as a starting point for your own `accounts.ini`
configuration file. This example configuration file attempts to illustrate
working values.

The current design (based off of the existing
<https://github.com/atc0005/list-emails> project) limits all email account
entries (reflected by different sections) to the same IMAP server. If you need
to process accounts from different servers you will need a separate copy of
the `accounts.ini` file for each server.

GH-122 is intended to change this in order to allow specifying the server/port
values per account entry. While a little redundant for accounts sharing the
same IMAP server, this is flexible enough to support separate account sections
for any number of IMAP servers.

Once reviewed and adjusted, your copy of the `accounts.ini` file can be placed
in one of the following locations to be automatically detected and used by
this application:

- alongside the `list-emails` (or `list-emails.exe`) binary (as
  `accounts.ini`)
- at `$HOME/.config/check-mail/accounts.ini` on a UNIX-like system (e.g.,
  Linux distro, Mac)
- at `C:\Users\YOUR_USERNAME\AppData\check-mail\accounts.ini` on a Windows
  system

You may also place the file wherever you like and refer to it using the
`-config-file` (full-length flag name). See the [Examples](#examples) and
[Command-line arguments](#command-line-arguments) sections for usage details.

## Examples

### `check_imap_mailbox`

#### As a Nagios plugin

When called by Nagios, you don't really benefit from having the application
generate log output; Nagios throws away output `stderr` and returns anything
sent to `stdout`, so output of any kind has to be carefully tailored to just
what you want to show up in the actual alert. Because of that, we disable
logging output explicitly and rely on the plugin to return information as
required via `stdout`.

```ShellSession
$ /usr/lib/nagios/plugins/check_imap_mailbox -folders "Inbox, Junk Email" -server imap.example.com -username "tacotuesdays@example.com" -port 993 -password "coconuts" -log-level disabled
OK: tacotuesdays@example.com: No messages found in folders: Inbox, Junk Email
```

#### Login failure

Assuming that an error occurred, we will want to explicitly choose a different
log level than the one normally used when the plugin is operating normally.
Here we choose `-log-level info` to get at basic operational details. You may
wish to use `-log-level debug` to get even more feedback.

```ShellSession
$ /usr/lib/nagios/plugins/check_imap_mailbox -folders "Inbox, Junk Email" -server imap.example.com -username "tacotuesdays@example.com" -port 993 -password "coconuts" -log-level info -branding
{"level":"error","username":"tacotuesdays@example.com","server":"imap.example.com","port":993,"folders_to_check":"Inbox,Junk Email","error":"LOGIN failed.","caller":"T:/github/check-mail/main.go:152","message":"Login error occurred"}
Login error occurred

Additional details: LOGIN failed.

Notification generated by check_imap_mailbox x.y.z
```

#### Help Output

```ShellSession
$ ./check_imap_mailbox --help
check-mail x.y.z
https://github.com/atc0005/check-mail

Usage of ./check_imap_mailbox:
  -branding
        Toggles emission of branding details with plugin status details. This output is disabled by default.
  -folders value
        Folders or IMAP "mailboxes" to check for mail. This value is provided as a comma-separated list.
  -log-level string
        Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace. (default "info")
  -password string
        The remote mail server account password.
  -port int
        TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections. (default 993)
  -server string
        The fully-qualified domain name of the remote mail server.
  -username string
        The account used to login to the remote mail server. This is often in the form of an email address.
  -version
        Whether to display application version and then immediately exit application.
```

### `list-emails`

#### No options

In this example, the `list-emails` application is in the current working
directory, as is the `accounts.ini` file. When run, the `output` and `log`
directories are created (if not already present) and populated with new log
and report files.

```ShellSession
$ ./list-emails
Checking account: email1
Checking account: email2
OK: Successfully generated reports for accounts: email1@example.com, email2@example.com
```

#### Alternate locations for config file, log and report directories

For this example, I intentionally placed each item on a separate volume. I
then reference each item via separate flags.

```ShellSession
./list-emails --config-file /mnt/t/accounts.ini --report-file-dir /mnt/g/reports --log-file-dir /mnt/d/log
Checking account: email1
Checking account: email2
OK: Successfully generated reports for accounts: email1@example.com, email2@example.com
```

#### Help output

```ShellSession
$ ./list-emails --help
check-mail x.y.z
https://github.com/atc0005/check-mail

Usage of ./list-emails:
  -config-file string
        Full path to the INI-formatted configuration file used by this application. See contrib/list-emails/accounts.example.ini for a starter template. Rename to accounts.ini, update with applicable information and place in a directory of your choice. If this file is found in your current working directory you need not use this flag. (default "accounts.ini")
  -log-file-dir string
        Full path to the directory where log files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist. (default "log")
  -log-level string
        Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace. (default "info")
  -report-file-dir string
        Full path to the directory where email summary report files will be created. The user account running this application requires write permission to this directory. If not specified, a default directory will be created in your current working directory if it does not already exist. (default "output")
  -version
        Whether to display application version and then immediately exit application.
```

## License

From the [LICENSE](LICENSE) file:

```license
MIT License

Copyright (c) 2020 Adam Chalkley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## References

- <https://github.com/emersion/go-imap>
- <https://github.com/rs/zerolog>
- <https://github.com/go-ini/ini>
- <https://github.com/atc0005/go-nagios>
