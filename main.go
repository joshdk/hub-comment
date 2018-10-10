// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const tokenEnvVar = "GITHUB_TOKEN"

var version = "development"

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "hub-comment: %s\n", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	token, found := os.LookupEnv(tokenEnvVar)
	if !found {
		return fmt.Errorf("no GITHUB_TOKEN set in environment")
	}

	var (
		ctx    = context.Background()
		client = makeClient(ctx, token)
	)

	// Get the current user associated with the given API token.
	self, err := getSelf(ctx, client)
	if err != nil {
		return err
	}

	fmt.Printf("User is %s (%s)\n", self.GetName(), self.GetLogin())

	return nil
}

// makeClient builds a GitHub client that is authenticated with the given token.
func makeClient(ctx context.Context, token string) *github.Client {
	var (
		tokenSource = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(ctx, tokenSource)
	)
	return github.NewClient(httpClient)
}

// getSelf retrieves the current authenticated user.
func getSelf(ctx context.Context, client *github.Client) (*github.User, error) {
	// A literal user value of "" retrieves the current user.
	user, _, err := client.Users.Get(ctx, "")
	return user, err
}
