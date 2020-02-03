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

USAGE

    $ ./check_imap_mailbox
    check_imap_mailbox x.y.z
    https://github.com/atc0005/check-mail

    Usage of ./check_imap_mailbox:
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


*/
package main
