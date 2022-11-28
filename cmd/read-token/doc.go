// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Small CLI app used to read an OAuth2 Client Credentials token from a file
// for use within a shell script. A separate tool is used to retrieve the
// token from an authority (e.g., via a cron job) and cache it for this tool
// to read.
//
// See our [GitHub repo]:
//
//   - to review documentation (including examples)
//   - for the latest code
//   - to file an issue or submit improvements for review and potential
//     inclusion into the project
//
// [GitHub repo]: https://github.com/atc0005/check-mail
package main
