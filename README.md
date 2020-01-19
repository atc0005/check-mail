# Check Mail

Nagios plugins used to monitor mail services

[![Latest Release](https://img.shields.io/github/release/atc0005/check-mail.svg?style=flat-square)](https://github.com/atc0005/check-mail/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/check-mail?status.svg)](https://godoc.org/github.com/atc0005/check-mail)
![Validate Codebase](https://github.com/atc0005/check-mail/workflows/Validate%20Codebase/badge.svg)
![Validate Docs](https://github.com/atc0005/check-mail/workflows/Validate%20Docs/badge.svg)

## Project home

See [our GitHub repo](https://github.com/atc0005/check-mail) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

This repo contains various Nagios plugins used to monitor mail services.

PLACEHOLDER - Add table of plugins here

- `check_imap_folder`

## Features

- PLACEHOLDER
- Go modules (vs classic `GOPATH` setup)

## Requirements

- Go 1.12+ (for building)
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

## How to install it

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
   - for current operating system with default `go` build options
     - `go build`
   - for all supported platforms
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Copy the applicable binary to whatever systems that need to run it
   1. Linux: `/tmp/check-mail/check-mail`
   1. Windows: `/tmp/check-mail/check-mail.exe`

## Configuration Options

### Command-line Arguments

| Option          | Required | Default        | Repeat | Possible                            | Description                                                                                                                                  |
| --------------- | -------- | -------------- | ------ | ----------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       | `false`        | No     | `0+`                                | Keep specified number of matching files.                                                                                                     |
| `console`       | No       | `false`        | No     | `true`, `false`                     | Dump CSV file equivalent to console.                                                                                                         |
| `csvfile`       | Yes      | *empty string* | No     | *valid file name characters*        | The fully-qualified path to a CSV file that this application should generate.                                                                |
| `excelfile`     | No       | *empty string* | No     | *valid file name characters*        | The fully-qualified path to a Microsoft Excel file that this application should generate.                                                    |
| `size`          | No       | `1` (byte)     | No     | `0+`                                | File size limit for evaluation. Files smaller than this will be skipped.                                                                     |
| `duplicates`    | No       | `2`            | No     | `2+`                                | Number of files of the same file size needed before duplicate validation logic is applied.                                                   |
| `ignore-errors` | No       | `false`        | No     | `true`, `false`                     | Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files. |
| `path`          | Yes      | *empty string* | Yes    | *one or more valid directory paths* | Path to process. This flag may be repeated for each additional path to evaluate.                                                             |
| `recurse`       | No       | `false`        | No     | `true`, `false`                     | Perform recursive search into subdirectories per provided path.                                                                              |

## Examples

PLACEHOLDER

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
