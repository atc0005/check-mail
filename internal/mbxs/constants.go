// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package mbxs

import "github.com/atc0005/check-mail/internal/textutils"

// DefaultReplacementString is used when replacing Unicode characters
// incompatible with a target character set. The common use case is
// substituting Unicode characters incompatible with the utf8mb3 character
// set.
const DefaultReplacementString string = textutils.EmojiScissors
