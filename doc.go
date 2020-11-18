/*

This repo contains various tools used to monitor mail services.

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/check-mail) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

PURPOSE

Monitor remote mail services.

FEATURES

• Nagios plugin for monitoring one or mail remote IMAP mailboxes

• Command-line tool for generating reports of specified folders/mailboxes

USAGE

    $ ./check_imap_mailbox
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

*/
package main
