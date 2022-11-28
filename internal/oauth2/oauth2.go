// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package oauth2

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// GetClientCredentialsToken receives OAuth2 Client Credentials / application
// registration details used to request a token from an authorization server
// and returns a new token or an error if one occurs.
func GetClientCredentialsToken(
	ctx context.Context,
	clientID string,
	clientSecret string,
	scopes []string,
	tokenEndpointURL string,
	maxAttempts int,
) (*oauth2.Token, error) {

	oauth2Config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenEndpointURL,
		Scopes:       scopes,
	}

	var token *oauth2.Token
	var result error

	// Attempt to retrieve token, retry up to maximum before giving up.
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		token, result = oauth2Config.Token(ctx)

		switch {

		// Error encountered. Wait briefly before trying again.
		case result != nil:
			time.Sleep(1 * time.Second)

		// Token validity failed (for reasons unknown). Wait briefly before
		// trying again.
		case !token.Valid():
			time.Sleep(1 * time.Second)

		default:
			// Successful retrieval, return token.
			return token, nil
		}

	}

	return nil, fmt.Errorf(
		"failed to retrieve token: %w",
		result,
	)

}
