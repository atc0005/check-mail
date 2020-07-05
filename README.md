# Check Mail

Various tools used to monitor mail services

[![Latest Release](https://img.shields.io/github/release/atc0005/check-mail.svg?style=flat-square)](https://github.com/atc0005/check-mail/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/check-mail?status.svg)](https://godoc.org/github.com/atc0005/check-mail)
![Validate Codebase](https://github.com/atc0005/check-mail/workflows/Validate%20Codebase/badge.svg)
![Validate Docs](https://github.com/atc0005/check-mail/workflows/Validate%20Docs/badge.svg)

- [Check Mail](#check-mail)
  - [Project home](#project-home)
  - [Overview](#overview)
  - [Features](#features)
    - [`check_imap_mailbox`](#check_imap_mailbox)
    - [`list-emails`](#list-emails)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Configuration Options](#configuration-options)
    - [Command-line Arguments](#command-line-arguments)
  - [Examples](#examples)
    - [As a Nagios plugin](#as-a-nagios-plugin)
    - [Output](#output)
      - [Login failure](#login-failure)
      - [Help Output](#help-output)
  - [License](#license)
  - [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/check-mail) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

This repo contains various tools used to monitor mail services.

| Tool Name            | Status   | Description                                                             |
| -------------------- | -------- | ----------------------------------------------------------------------- |
| `check_imap_mailbox` | Alpha    | Nagios plugin used to monitor mailboxes for items                       |
| `list-emails`        | Planning | PLACEHOLDER: Small CLI app used to generate listing of mailbox contents |

## Features

### `check_imap_mailbox`

- Monitor one or many mailboxes
- Leveled logging
  - JSON-format output
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`.
- TLS/SSL IMAP4 connectivity (defaults to 993/tcp)
- Go modules (vs classic `GOPATH` setup)

### `list-emails`

Placeholder. See GH-2 for details.

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
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Copy the applicable binary to whatever systems needs to run it
   - if using `Makefile`: look in `/tmp/check-mail/release_assets/check_imap_mailbox/`
   - if using `go build`: look in `/tmp/check-mail/`

## Configuration Options

### Command-line Arguments

| Option          | Required | Default        | Repeat | Possible                                                                | Description                                                                                                                                                                                 |
| --------------- | -------- | -------------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       |                | No     | `-h`, `--help`                                                          | Generate listing of all valid command-line options and applicable (short) guidance for using them.                                                                                          |
| `folders`       | Yes      | *empty string* | No     | *comma-separated list of folders*                                       | Folders or IMAP "mailboxes" to check for mail. This value is provided as a comma-separated list.                                                                                            |
| `username`      | Yes      | *empty string* | No     | *valid username, often in email address format*                         | The account used to login to the remote mail server. This is often in the form of an email address.                                                                                         |
| `password`      | Yes      | *empty string* | No     | *valid password*                                                        | The remote mail server account password.                                                                                                                                                    |
| `server`        | Yes      | *empty string* | No     | *valid FQDN or IP Address*                                              | The fully-qualified domain name of the remote mail server.                                                                                                                                  |
| `port`          | No       | `993`          | No     | *valid IMAP TCP port*                                                   | TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections.                                                                  |
| `logging-level` | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                                                                                             |
| `branding`      | No       | `false`        | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. Because this output may not mix well with branding information emitted by other tools, this output is disabled by default. |

## Examples

### As a Nagios plugin

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

### Output

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
check_imap_mailbox x.y.z
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
- <https://github.com/atc0005/go-nagios>
