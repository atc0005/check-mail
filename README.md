<!-- omit in toc -->
# Check Mail

Various tools used to monitor mail services

[![Latest Release](https://img.shields.io/github/release/atc0005/check-mail.svg?style=flat-square)](https://github.com/atc0005/check-mail/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-mail.svg)](https://pkg.go.dev/github.com/atc0005/check-mail)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/check-mail)](https://github.com/atc0005/check-mail)
[![Lint and Build](https://github.com/atc0005/check-mail/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/check-mail/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/check-mail/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/check-mail/actions/workflows/project-analysis.yml)
[![Push Validation](https://github.com/atc0005/check-mail/actions/workflows/push-validation.yml/badge.svg)](https://github.com/atc0005/check-mail/actions/workflows/push-validation.yml)

<!-- omit in toc -->
## Table of contents

- [Project home](#project-home)
- [Overview](#overview)
- [Features](#features)
  - [`check_imap_mailbox_*`](#check_imap_mailbox_)
  - [`list-emails`](#list-emails)
  - [`lsimap`](#lsimap)
- [Requirements](#requirements)
  - [Building source code](#building-source-code)
  - [Running](#running)
  - [Office 365 (O365) permissions](#office-365-o365-permissions)
- [Installation](#installation)
  - [From source](#from-source)
  - [Using release binaries](#using-release-binaries)
- [Configuration Options](#configuration-options)
  - [`check_imap_mailbox_basic`](#check_imap_mailbox_basic)
    - [Command-line arguments](#command-line-arguments)
  - [`check_imap_mailbox_oauth2`](#check_imap_mailbox_oauth2)
    - [Required preparation](#required-preparation)
    - [Command-line arguments](#command-line-arguments-1)
  - [`list-emails`](#list-emails-1)
    - [Command-line arguments](#command-line-arguments-2)
    - [Configuration file](#configuration-file)
      - [Settings](#settings)
        - [Basic Auth](#basic-auth)
        - [OAuth2](#oauth2)
      - [Usage](#usage)
  - [`lsimap`](#lsimap-1)
    - [Command-line arguments](#command-line-arguments-3)
  - [`xoauth2`](#xoauth2)
    - [Command-line arguments](#command-line-arguments-4)
- [Examples](#examples)
  - [`check_imap_mailbox_basic`](#check_imap_mailbox_basic-1)
    - [As a Nagios plugin](#as-a-nagios-plugin)
    - [Login failure](#login-failure)
  - [`check_imap_mailbox_oauth2`](#check_imap_mailbox_oauth2-1)
  - [`list-emails`](#list-emails-2)
    - [No options](#no-options)
    - [Alternate locations for config file, log and report directories](#alternate-locations-for-config-file-log-and-report-directories)
  - [`lsimap`](#lsimap-2)
- [License](#license)
- [References](#references)
  - [Related projects](#related-projects)
  - [Dependencies](#dependencies)
  - [OAuth2 Research](#oauth2-research)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Overview

This repo contains various tools used to monitor mail services.

| Tool Name                   | Overall Status | Description                                                                                                      |
| --------------------------- | -------------- | ---------------------------------------------------------------------------------------------------------------- |
| `check_imap_mailbox_basic`  | Stable         | Nagios plugin used to monitor mailboxes for items (via Basic Auth)                                               |
| `check_imap_mailbox_oauth2` | Alpha          | Nagios plugin used to monitor mailboxes for items (via OAuth2)                                                   |
| `list-emails`               | Stable         | Small CLI app used to generate listing of mailbox contents                                                       |
| `lsimap`                    | Alpha          | Small CLI tool to list advertised capabilities for specified IMAP server                                         |
| `xoauth2`                   | Alpha          | Small CLI tool to convert given username and token to XOAuth2 formatted (optionally SASL XOAUTH2 encoded) string |

## Features

### `check_imap_mailbox_*`

There are two plugins which perform the same overall function, but utilize
different mechanisms to authenticate to a specific IMAP server:

- `check_imap_mailbox_basic`
  - uses Basic Auth (username/password) for authentication
- `check_imap_mailbox_oauth2`
  - uses OAuth2 Client Credentials (client ID/secret) flow for authentication

Shared functionality:

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
  - user-specified minimum TLS version
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
- Multiple authentication options
  - Basic Auth (username/password)
  - OAuth2 Client Credentials (client ID/secret) flow
- TLS IMAP4 connectivity
  - port defaults to 993/tcp
  - network type defaults to either of IPv4 and IPv6, but optionally limited
    to IPv4-only or IPv6-only
  - user-specified minimum TLS version
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

### `lsimap`

- Quick one-off tool to list advertised capabilities for specified IMAP server
- Leveled logging
  - `console writer`: human-friendly, colorized output
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`
  - enable `debug` level to monitor submitted IMAP commands and received IMAP
    server responses
- TLS IMAP4 connectivity
  - port defaults to 993/tcp
  - network type defaults to either of IPv4 and IPv6, but optionally limited
    to IPv4-only or IPv6-only
  - user-specified minimum TLS version

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for *preferred* version
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`
- `tc-hib/go-winres`
  - if using the provided `Makefile`
  - used to generate Windows resource files

### Running

- Windows 10
- Ubuntu Linux 18.04+

### Office 365 (O365) permissions

The `list-emails` and `check_imap_mailbox_oauth2` tools support OAuth2 Client
Credentials flow authentication.

The [`check_imap_mailbox_oauth2`](#check_imap_mailbox_oauth2-1) example
illustrates connecting to an Office 365 (O365) "shared mailbox" (aka, a shared
account).

- The `client-id`, `client-secret` flag values are obtained from the
  [application registration][azure-app-registration].
- The `https://outlook.office365.com/.default` scopes value indicates that the
  permissions listed in the application registration should be used.

Testing was performed with these permissions/scopes set within the application
registration:

| API                          | Permissions name        | Type        | Description                                  | Admin consent required |
| ---------------------------- | ----------------------- | ----------- | -------------------------------------------- | ---------------------- |
| `Microsoft Graph`            | `IMAP.AccessAsUser.All` | Delegated   | Read and write access to mailboxes via IMAP. | No                     |
| `Microsoft Graph`            | `User.Read`             | Delegated   | Sign in and read user profile                | No                     |
| `Office 365 Exchange Online` | `IMAP.AccessAsApp`      | Application | `IMAP.AccessAsApp`                           | Yes                    |

The last one has to be granted by a tenant administrator.

Per <https://blog.rebex.net/office365-ews-oauth-unattended>:

> Optionally, you can remove the delegated `User.Read` permission which is not
> needed for app-only application - click the context menu on the right side
> of the permission and select Remove permission.

Other sources have said the same thing: Delegated scopes are not needed for
the `client credentials` flow; only the `IMAP.AccessAsApp` permission is
required for the OAuth2 Client Credentials flow (used by tools in this
project).

Lastly, an Office 365 tenant administrator needs to:

1. [register the service principals in Exchange][azure-app-register-service-principals]
1. add specific mailboxes in the tenant that will be allowed to be accessed by this plugin

See the [`check_imap_mailbox_oauth2`](#check_imap_mailbox_oauth2-1) example or
the official O365 [o365-cred-flow-test-script] test script to confirm that
required settings are in place.

## Installation

### From source

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
   - NOTE: Pay special attention to the remarks about `$HOME/.profile`
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
   - for current operating system (using bundled dependencies)
     - `go build -mod=vendor ./cmd/check_imap_mailbox_basic/`
     - `go build -mod=vendor ./cmd/check_imap_mailbox_oauth2/`
     - `go build -mod=vendor ./cmd/list-emails/`
     - `go build -mod=vendor ./cmd/lsimap/`
     - `go build -mod=vendor ./cmd/xoauth2/`
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Locate generated binaries
   - if using `Makefile`
     - look in `/tmp/check-mail/release_assets/check_imap_mailbox_basic/`
     - look in `/tmp/check-mail/release_assets/check_imap_mailbox_oauth2/`
     - look in `/tmp/check-mail/release_assets/list-emails/`
     - look in `/tmp/check-mail/release_assets/lsimap/`
     - look in `/tmp/check-mail/release_assets/xoauth2/`
   - if using `go build`
     - look in `/tmp/check-mail/`
1. Copy the applicable binaries to whatever systems needs to run them
1. Deploy
   - Place `list-emails` in a location of your choice
   - Place `lsimap` in a location of your choice
   - Place `xoauth2` in a location of your choice
   - Place `check_imap_mailbox_basic` in the same location where your distro's
     package manage has place other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_imap_mailbox_basic` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_imap_mailbox_basic` on RedHat-based
       systems
   - Place `check_imap_mailbox_oauth2` in the same location where your distro's
     package manage has place other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_imap_mailbox_oauth2` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_imap_mailbox_oauth2` on RedHat-based
       systems
1. Copy the template [configuration file](#configuration-file), modify
   accordingly and place in a [supported location](#configuration-file)

### Using release binaries

1. Download the [latest
   release](https://github.com/atc0005/check-mail/releases/latest) binaries
1. Deploy
   - Place `list-emails` in a location of your choice
   - Place `lsimap` in a location of your choice
   - Place `xoauth2` in a location of your choice
   - Place `check_imap_mailbox_basic` in the same location where your distro's
     package manager places other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_imap_mailbox_basic` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_imap_mailbox_basic` on RedHat-based
       systems
   - Place `check_imap_mailbox_oauth2` in the same location where your distro's
     package manager places other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_imap_mailbox_oauth2` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_imap_mailbox_oauth2` on RedHat-based
       systems
1. Copy the template [configuration file](#configuration-file), modify
   accordingly and place in a [supported location](#configuration-file)

## Configuration Options

### `check_imap_mailbox_basic`

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
| `min-tls`       | No       | `tls12`        | No     | `tls10`, `tls11`, `tls12`, `tls13`                                      | Limits version of TLS used for connections to remote mail servers.                                                                                                                          |
| `logging-level` | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                                                                                             |
| `branding`      | No       | `false`        | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. Because this output may not mix well with branding information emitted by other tools, this output is disabled by default. |
| `version`       | No       | `false`        | No     | `true`, `false`                                                         | Whether to display application version and then immediately exit application                                                                                                                |

### `check_imap_mailbox_oauth2`

#### Required preparation

This plugin uses the OAuth2 Client Credentials flow to authenticate.

This requires registering an application with the authority for the resource
that you wish to access. The specifics differ (at least slightly) for every
IMAP account provider that you wish to interact with.

See the [Office 365 (O365) permissions](#office-365-o365-permissions) section
for details specific to using this plugin with O365 mailboxes.

#### Command-line arguments

- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined.

| Option           | Required | Default        | Repeat | Possible                                                                | Description                                                                                                                                                                                     |
| ---------------- | -------- | -------------- | ------ | ----------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`      | No       |                | No     | `-h`, `--help`                                                          | Generate listing of all valid command-line options and applicable (short) guidance for using them.                                                                                              |
| `folders`        | Yes      | *empty string* | No     | *comma-separated list of folders*                                       | Folders or IMAP "mailboxes" to check for mail. This value is provided as a comma-separated list.                                                                                                |
| `scopes`         | Yes      | *empty string* | No     | *comma-separated list of scopes*                                        | Permissions needed by the application. If using the scopes defined by the application registration you must use the `RESOURCE/.default` format (e.g., `https://outlook.office365.com/.default`. |
| `client-id`      | Yes      | *empty string* | No     | *valid application ID associated with registered app*                   | Application (client) ID created during app registration.                                                                                                                                        |
| `client-secret`  | Yes      | *empty string* | No     | *valid application secret associated with registered app*               | Client secret (aka, "app" password).                                                                                                                                                            |
| `shared-mailbox` | Yes      | *empty string* | No     | *valid shared mailbox name, often in email address format*              | Email account that is to be accessed using client ID & secret values. Usually a shared mailbox among a team.                                                                                    |
| `port`           | No       | `993`          | No     | *valid IMAP TCP port*                                                   | TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections.                                                                      |
| `net-type`       | No       | `auto`         | No     | `auto`, `tcp4`, `tcp6`                                                  | Limits network connections to remote mail servers to one of the specified types.                                                                                                                |
| `min-tls`        | No       | `tls12`        | No     | `tls10`, `tls11`, `tls12`, `tls13`                                      | Limits version of TLS used for connections to remote mail servers.                                                                                                                              |
| `logging-level`  | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                                                                                                 |
| `branding`       | No       | `false`        | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. Because this output may not mix well with branding information emitted by other tools, this output is disabled by default.     |
| `version`        | No       | `false`        | No     | `true`, `false`                                                         | Whether to display application version and then immediately exit application                                                                                                                    |

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
| `min-tls`         | No       | `tls12`        | No     | `tls10`, `tls11`, `tls12`, `tls13`                                      | Limits version of TLS used for connections to remote mail servers.                                                                                                                                                                                                                                                                       |
| `logging-level`   | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                                                                                                                                                                                                                                          |
| `version`         | No       | `false`        | No     | `true`, `false`                                                         | Whether to display application version and then immediately exit application                                                                                                                                                                                                                                                             |

#### Configuration file

##### Settings

**NOTE**: The `email1` and `email2` value below is for illustration. You are
free to choose section names, though it is recommended to base them off of the
username (sans `@` symbol and domain part) for each email account. While only
`email1` is listed, many such entries (one per account) are supported.

The `list-emails` CLI app supports both Basic Auth and OAuth2 Client
Credentials flow for authentication. Depending on the desired authentication
type some settings are required, others ignored; if using Basic Auth settings
specific to OAuth2 are ignored.

###### Basic Auth

| Config file Setting Name | Section Name | Notes                                               |
| ------------------------ | ------------ | --------------------------------------------------- |
| `server_name`            | `DEFAULT`    | FQDN of IMAP server (e.g., `outlook.office365.com`) |
| `server_port`            | `DEFAULT`    | Usually 993                                         |
| `username`               | `email1`     | Often in the form of an email address               |
| `password`               | `email1`     | Account password                                    |
| `folders`                | `email1`     | Double quoted, comma separated                      |

###### OAuth2

| Config file Setting Name | Section Name | Notes                                                                                                          |
| ------------------------ | ------------ | -------------------------------------------------------------------------------------------------------------- |
| `server_name`            | `DEFAULT`    | FQDN of IMAP server (e.g., `outlook.office365.com`)                                                            |
| `server_port`            | `DEFAULT`    | Usually 993                                                                                                    |
| `client_id`              | `DEFAULT`    | The ID associated with the application registration                                                            |
| `client_secret`          | `DEFAULT`    | Application secret (aka, "app" password)                                                                       |
| `scopes`                 | `DEFAULT`    | Comma-separated list of permissions needed by the application (e.g., `https://outlook.office365.com/.default`) |
| `endpoint_token_url`     | `DEFAULT`    | The OAuth2 provider's token endpoint URL.                                                                      |
| `shared_mailbox`         | `email1`     | Email address format (e.g., `me@there.com`)                                                                    |
| `folders`                | `email1`     | Double quoted, comma separated                                                                                 |

##### Usage

There are two example INI files available which illustrate available
configuration settings:

- [contrib/list-emails/basic-auth/accounts.example.ini](contrib/list-emails/basic-auth/accounts.example.ini)
  - Basic Auth (username/password) for authentication
- [contrib/list-emails/oauth2/accounts.example.ini](contrib/list-emails/oauth2/accounts.example.ini)
  - OAuth2 Client Credentials (client ID/secret) flow for authentication

These files are intended as starting points for your own `accounts.ini`
configuration file.

The current design (based off of the existing
<https://github.com/atc0005/list-emails> project) limits all email account
entries (reflected by different sections) to the same IMAP server. If you need
to process accounts from different servers you will need a separate copy of
the `accounts.ini` file for each server.

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

### `lsimap`

#### Command-line arguments

- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined.

| Option          | Required | Default        | Repeat | Possible                                                                | Description                                                                                                                |
| --------------- | -------- | -------------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       |                | No     | `-h`, `--help`                                                          | Generate listing of all valid command-line options and applicable (short) guidance for using them.                         |
| `server`        | Yes      | *empty string* | No     | *valid FQDN or IP Address*                                              | The fully-qualified domain name of the remote mail server.                                                                 |
| `port`          | No       | `993`          | No     | *valid IMAP TCP port*                                                   | TCP port used to connect to the remote mail server. This is usually the same port used for TLS encrypted IMAP connections. |
| `net-type`      | No       | `auto`         | No     | `auto`, `tcp4`, `tcp6`                                                  | Limits network connections to remote mail servers to one of the specified types.                                           |
| `min-tls`       | No       | `tls12`        | No     | `tls10`, `tls11`, `tls12`, `tls13`                                      | Limits version of TLS used for connections to remote mail servers.                                                         |
| `logging-level` | No       | `info`         | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Sets log level.                                                                                                            |
| `version`       | No       | `false`        | No     | `true`, `false`                                                         | Whether to display application version and then immediately exit application                                               |

### `xoauth2`

#### Command-line arguments

- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined.

| Option      | Required | Default        | Repeat | Possible             | Description                                                                                        |
| ----------- | -------- | -------------- | ------ | -------------------- | -------------------------------------------------------------------------------------------------- |
| `h`, `help` | No       |                | No     | `-h`, `--help`       | Generate listing of all valid command-line options and applicable (short) guidance for using them. |
| `account`   | Yes      | *empty string* | No     | *valid account name* | Username or mailbox in email format.                                                               |
| `token`     | Yes      | *empty string* | No     | *valid token*        | Access token.                                                                                      |
| `encode`    | No       | `false`        | No     | `true`, `false`      | Whether to encode XOAuth2 string for use in SASL XOAUTH2.                                          |

## Examples

### `check_imap_mailbox_basic`

#### As a Nagios plugin

When called by Nagios, you don't really benefit from having the application
generate log output; Nagios throws away output `stderr` and returns anything
sent to `stdout`, so output of any kind has to be carefully tailored to just
what you want to show up in the actual alert. Because of that, we disable
logging output explicitly and rely on the plugin to return information as
required via `stdout`.

```ShellSession
$ /usr/lib/nagios/plugins/check_imap_mailbox_basic -folders "Inbox, Junk Email" -server imap.example.com -username "tacotuesdays@example.com" -port 993 -password "coconuts" -log-level disabled
OK: tacotuesdays@example.com: No messages found in folders: Inbox, Junk Email
```

#### Login failure

Assuming that an error occurred, we will want to explicitly choose a different
log level than the one normally used when the plugin is operating normally.
Here we choose `-log-level info` to get at basic operational details. You may
wish to use `-log-level debug` to get even more feedback.

```ShellSession
$ /usr/lib/nagios/plugins/check_imap_mailbox_basic -folders "Inbox, Junk Email" -server imap.example.com -username "tacotuesdays@example.com" -port 993 -password "coconuts" -log-level info -branding
{"level":"error","username":"tacotuesdays@example.com","server":"imap.example.com","port":993,"folders_to_check":"Inbox,Junk Email","error":"LOGIN failed.","caller":"T:/github/check-mail/main.go:152","message":"Login error occurred"}
Login error occurred

Additional details: LOGIN failed.

Notification generated by check_imap_mailbox_basic x.y.z
```

### `check_imap_mailbox_oauth2`

Aside from accepting a different set of flags and authenticating using OAuth2
Client Credentials flow, the functionality of this plugin is identical to
`check_imap_mailbox_basic`.

```ShellSession
$ /usr/lib/nagios/plugins/check_imap_mailbox_basic --shared-mailbox "tacotuesdays@example.com" --folders "Inbox, Junk Email" --server outlook.office365.com --client-id "ZYDPLLBWSK3MVQJSIYHB1OR2JXCY0X2C5UJ2QAR2MAAIT5Q" --client-secret "_djgA8heFo0WSIMom7U39WmGTQFHWkcD8x-A1o-4sro" --token-url "https://login.microsoftonline.com/6029c1d9-aa2f-4227-8f7c-0c23224a0fa9/oauth2/v2.0/token" --scopes "https://outlook.office365.com/.default" --port 993 --log-level disabled
OK: tacotuesdays@example.com: No messages found in folders: Inbox, Junk Email
```

See the [Office 365 (O365) permissions](#office-365-o365-permissions) section
for details specific to using this plugin with O365 mailboxes.

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
OK: Successfully generated reports for accounts: email1, email2
```

#### Alternate locations for config file, log and report directories

For this example, I intentionally placed each item on a separate volume. I
then reference each item via separate flags.

```ShellSession
./list-emails --config-file /mnt/t/accounts.ini --report-file-dir /mnt/g/reports --log-file-dir /mnt/d/log
Checking account: email1
Checking account: email2
OK: Successfully generated reports for accounts: email1, email2
```

### `lsimap`

Quick listings for outlook.office365.com and imap.gmail.com.

This tool can be useful for determining at a glance what authentication
mechanisms are supported by an IMAP server.

```console
$ ./lsimap --server outlook.office365.com
6:10AM INF cmd\lsimap\main.go:61 > Connection established to server
6:10AM INF cmd\lsimap\main.go:70 > Gathering pre-login capabilities
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=PLAIN
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=XOAUTH2
6:10AM INF cmd\lsimap\main.go:87 > Capability: CHILDREN
6:10AM INF cmd\lsimap\main.go:87 > Capability: ID
6:10AM INF cmd\lsimap\main.go:87 > Capability: IDLE
6:10AM INF cmd\lsimap\main.go:87 > Capability: IMAP4
6:10AM INF cmd\lsimap\main.go:87 > Capability: IMAP4rev1
6:10AM INF cmd\lsimap\main.go:87 > Capability: LITERAL+
6:10AM INF cmd\lsimap\main.go:87 > Capability: MOVE
6:10AM INF cmd\lsimap\main.go:87 > Capability: NAMESPACE
6:10AM INF cmd\lsimap\main.go:87 > Capability: SASL-IR
6:10AM INF cmd\lsimap\main.go:87 > Capability: UIDPLUS
6:10AM INF cmd\lsimap\main.go:87 > Capability: UNSELECT
6:10AM INF cmd\lsimap\main.go:95 > Connection to server closed

$ ./lsimap --server imap.gmail.com
6:10AM INF cmd\lsimap\main.go:61 > Connection established to server
6:10AM INF cmd\lsimap\main.go:70 > Gathering pre-login capabilities
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=OAUTHBEARER
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=PLAIN
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=PLAIN-CLIENTTOKEN
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=XOAUTH
6:10AM INF cmd\lsimap\main.go:87 > Capability: AUTH=XOAUTH2
6:10AM INF cmd\lsimap\main.go:87 > Capability: CHILDREN
6:10AM INF cmd\lsimap\main.go:87 > Capability: ID
6:10AM INF cmd\lsimap\main.go:87 > Capability: IDLE
6:10AM INF cmd\lsimap\main.go:87 > Capability: IMAP4rev1
6:10AM INF cmd\lsimap\main.go:87 > Capability: NAMESPACE
6:10AM INF cmd\lsimap\main.go:87 > Capability: QUOTA
6:10AM INF cmd\lsimap\main.go:87 > Capability: SASL-IR
6:10AM INF cmd\lsimap\main.go:87 > Capability: UNSELECT
6:10AM INF cmd\lsimap\main.go:87 > Capability: X-GM-EXT-1
6:10AM INF cmd\lsimap\main.go:87 > Capability: XLIST
6:10AM INF cmd\lsimap\main.go:87 > Capability: XYZZY
6:10AM INF cmd\lsimap\main.go:95 > Connection to server closed
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

### Related projects

- [Monitoring plugins](https://github.com/atc0005?tab=repositories&q=check-&type=source&language=go&sort=)
- <https://github.com/atc0005/send2teams>
- <https://github.com/atc0005/nagios-debug>

### Dependencies

- <https://github.com/emersion/go-imap>
- <https://github.com/emersion/go-sasl>
- <https://github.com/rs/zerolog>
- <https://github.com/go-ini/ini>
- <https://github.com/atc0005/go-nagios>
- <https://github.com/tc-hib/go-winres>
- <https://github.com/sqs/go-xoauth2>
- <https://golang.org/x/oauth2>

### OAuth2 Research

General:

- <https://github.com/atc0005/check-mail/issues/313>
  - my notes while researching/testing OAuth2 Client Credentials flow support
- <https://alexbilbie.com/guide-to-oauth-2-grants/>
- <https://aaronparecki.com/oauth-2-simplified/#authorization>
- <https://learn.microsoft.com/en-us/exchange/clients-and-mobile-in-exchange-online/deprecation-of-basic-authentication-exchange-online>
- <https://learn.microsoft.com/en-us/exchange/client-developer/legacy-protocols/how-to-authenticate-an-imap-pop-smtp-application-by-using-oauth>
- <https://learn.microsoft.com/en-us/azure/active-directory/develop/active-directory-v2-protocols>
- <https://learn.microsoft.com/en-us/azure/active-directory/develop/quickstart-register-app>
- <https://learn.microsoft.com/en-us/azure/active-directory/develop/scenario-daemon-app-registration>
- <https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-client-creds-grant-flow>
- <https://www.oauth.com/oauth2-servers/access-tokens/access-token-response/>
- <https://blog.rebex.net/office365-ews-oauth-unattended>
- <https://www.youtube.com/watch?v=bMYA-146dmM>
- <https://stackoverflow.com/questions/73463357/cannot-authenticate-to-imap-on-office365-using-javamail>
- <https://learn.microsoft.com/en-us/azure/active-directory/develop/app-objects-and-service-principals>

Redmine:

- <https://www.redmine.org/issues/37688>
- <https://ruby-doc.org/stdlib-2.7.0/libdoc/net/imap/rdoc/Net/IMAP.html>
- <https://github.com/Mailbutler/mail_xoauth2>
  - <https://github.com/Mailbutler/mail_xoauth2/blob/master/lib/mail_xoauth2/oauth2_string.rb>
  - <https://github.com/Mailbutler/mail_xoauth2/blob/master/lib/mail_xoauth2/imap_xoauth2_authenticator.rb>

OAuth 2.0 Resource Owner Password Credentials (ROPC) grant:

- <https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth-ropc>
  - not recommended
  - deprecated, removed in OAuth 2.1
  - only a matter of time before it is removed
- <https://oauth.net/2/grant-types/password/>
- <https://www.oauth.com/oauth2-servers/access-tokens/password-grant/>

Some Go-specific references:

- <https://pkg.go.dev/golang.org/x/oauth2>
- <https://pkg.go.dev/golang.org/x/oauth2/microsoft>
- <https://github.com/sqs/go-xoauth2>
  - <https://github.com/sqs/go-xoauth2/issues/1>
- <https://github.com/AzureAD/microsoft-authentication-library-for-go>
- <https://github.com/google/gmail-oauth2-tools/tree/master/go/sendgmail>
- <https://github.com/emersion/go-imap>
  - not OAuth2 specific, but used by tools in this repo
  - <https://github.com/emersion/go-imap/wiki/Using-authentication-mechanisms>
- <https://github.com/emersion/go-sasl>
  - <https://github.com/emersion/go-sasl/issues/18>
  - XOAUTH2 support previously supplied by this project, now bundled locally

RFCS:

- <https://datatracker.ietf.org/doc/html/rfc6749#section-4.3>
- <https://datatracker.ietf.org/doc/html/rfc6749#section-4.4.3>
- <https://datatracker.ietf.org/doc/html/rfc9051#section-2.2>

Other projects:

- <https://github.com/tgulacsi/imapclient>
- <https://github.com/DanijelkMSFT/ThisandThat/blob/main/Get-IMAPAccessToken.ps1>
  - useful test script to ensure that credential flow is functional for a registered app and associated accounts
  - <https://learn.microsoft.com/en-us/exchange/clients-and-mobile-in-exchange-online/deprecation-of-basic-authentication-exchange-online>
  - <https://www.linkedin.com/pulse/start-using-oauth-office-365-popimap-authentication-danijel-klaric>
  - <https://techcommunity.microsoft.com/t5/exchange-team-blog/basic-authentication-deprecation-in-exchange-online-may-2022/ba-p/3301866>
  - <https://techcommunity.microsoft.com/t5/exchange-team-blog/announcing-oauth-2-0-client-credentials-flow-support-for-pop-and/ba-p/3562963>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/send2teams>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"

[azure-app-registration]:
    <https://learn.microsoft.com/en-us/azure/active-directory/develop/scenario-daemon-app-registration>
    "Azure App Registration"

[azure-app-register-service-principals]:
    <https://learn.microsoft.com/en-us/exchange/client-developer/legacy-protocols/how-to-authenticate-an-imap-pop-smtp-application-by-using-oauth#register-service-principals-in-exchange>
    "Register service principals in Exchange"

[o365-cred-flow-test-script]: <https://github.com/DanijelkMSFT/ThisandThat/blob/main/Get-IMAPAccessToken.ps1> "O365 Client Credentials flow test script"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
