/*

Package nagios provides common types and package-level variables for use with
Nagios plugins.

OVERVIEW

This package contains common types and package-level variables. The intent is
to reduce code duplication between various plugins that we maintain.

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/go-nagios) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

HOW TO USE

Assuming that you're using Go Modules
(https://blog.golang.org/using-go-modules), add this line to your imports like
so:

    package main

    import (
    "fmt"
    "log"
    "os"

    "github.com/atc0005/go-nagios"
    )

Then in your code reference the data types as you would from any other
package:

    fmt.Println("OK: All checks have passed")
    os.Exit(nagios.StateOK)

When you next build your package this one should be pulled in.

*/
package nagios
