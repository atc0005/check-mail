// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atc0005/check-mail/internal/config"
	"github.com/sqs/go-xoauth2"
)

type credentials struct {
	account string
	token   string
}

// Usage is a custom override for the default Help text provided by the flag
// package. Here we prepend some additional metadata to the existing output.
func Usage() {
	fmt.Fprintln(flag.CommandLine.Output(), "\n"+config.Version()+"\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	creds := credentials{}
	var encode bool

	flag.StringVar(&creds.account, "account", "", "username or mailbox in email format")
	flag.StringVar(&creds.token, "token", "", "access token")
	flag.BoolVar(&encode, "encode", false, "encode XOAuth2 string for use in SASL XOAUTH2")
	flag.Usage = Usage
	flag.Parse()

	switch {
	case creds.account == "":
		fmt.Println("missing account!")
		fmt.Printf("\nUsage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	case creds.token == "":
		fmt.Println("missing token!")
		fmt.Printf("\nUsage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var xoauth2str string
	switch {
	case encode:
		// base64-encoded XOAUTH2 string, suitable for direct use in SASL XOAUTH2
		xoauth2str = xoauth2.XOAuth2String(creds.account, creds.token)
	default:
		// unencoded XOAUTH2 string
		xoauth2str = xoauth2.OAuth2String(creds.account, creds.token)
	}

	// Intentionally skip emitting newline so we don't affect any wrapper
	// script that captures the output.
	fmt.Print(xoauth2str)

}
