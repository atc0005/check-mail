<!-- omit in toc -->
# go-nagios

Shared Golang package for Nagios plugins

[![Latest Release](https://img.shields.io/github/release/atc0005/go-nagios.svg?style=flat-square)](https://github.com/atc0005/go-nagios/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/go-nagios?status.svg)](https://godoc.org/github.com/atc0005/go-nagios)
[![Validate Codebase](https://github.com/atc0005/go-nagios/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/go-nagios/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/go-nagios/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/go-nagios/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/go-nagios/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/go-nagios/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/go-nagios/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/go-nagios/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of contents

- [Overview](#overview)
- [Features](#features)
- [Changelog](#changelog)
- [How to use](#how-to-use)
- [License](#license)
- [References](#references)

## Overview

This package contains common types and package-level variables used when
developing Nagios plugins. The intent is to reduce code duplication between
various plugins and help reduce typos associated with literal strings.

## Features

- Nagios state constants

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## How to use

Assuming that you're using [Go
Modules](https://blog.golang.org/using-go-modules), add this line to your
imports like so:

```golang
package main

import (
  "fmt"
  "log"
  "os"

  "github.com/atc0005/go-nagios"
)
  ```

Then in your code, reference the data types as you would from any other
package:

```golang
fmt.Println("OK: All checks have passed")
os.Exit(nagios.StateOKExitCode)
```

Alternatively, you can also use the provided state "labels" (constants) to
avoid using literal string state values:

```golang
fmt.Printf("%s: All checks have passed\r\n", nagios.StateOKLabel)
os.Exit(nagios.StateOKExitCode)
```

When you next build your package this one should be pulled in.

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

- <https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins-scripts/utils.sh.in>
- <http://nagios-plugins.org/doc/guidelines.html>
