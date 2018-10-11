// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/joshdk/hub-comment/hub"
	"golang.org/x/oauth2"
)

const (
	// githubTokenEnvVar is the name of the environment variable which holds
	// the token used for authenticating against the GitHub API.
	githubTokenEnvVar = "GITHUB_TOKEN"

	// pullRequestLinkEnvVar is the name of the environment variable which
	// holds a link to the current pull request. Injected automatically by
	// CircleCI.
	pullRequestLinkEnvVar = "CIRCLE_PULL_REQUEST"
)

var version = "development"

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "hub-comment: %s\n", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	token, found := os.LookupEnv(githubTokenEnvVar)
	if !found {
		return fmt.Errorf("no GITHUB_TOKEN set in environment")
	}

	// If no CIRCLE_PULL_REQUEST is set, print an error and return immediately
	// but do not fail. This environment variable will not be set on non-pr
	// branches, or if a build is started before a pr is opened.
	reference, found := os.LookupEnv(pullRequestLinkEnvVar)
	if !found {
		fmt.Fprintln(os.Stderr, "hub-comment: no CIRCLE_PULL_REQUEST set in environment")
		return nil
	}

	owner, repo, number, found := hub.SplitPullRequestReference(reference)
	if !found {
		return fmt.Errorf("malformed pull request link")
	}

	if len(os.Args) != 2 {
		return fmt.Errorf("no comment given")
	}

	text := os.Args[1]

	var (
		ctx    = context.Background()
		client = makeClient(ctx, token)
	)

	// Get the current user associated with the given API token.
	self, err := getSelf(ctx, client)
	if err != nil {
		return err
	}

	// Get a list of all comments for the given PR number
	comments, err := hub.GetComments(ctx, client, owner, repo, number)
	if err != nil {
		return err
	}

	// Select the most recent comment that was authored by the current user, if
	// one exists.
	commentID, found := hub.FilterComments(comments, self.GetLogin())

	fmt.Printf("User is %s (%s)\n", self.GetName(), self.GetLogin())
	fmt.Printf("Pull is %s/%s #%d\n", owner, repo, number)
	if found {
		fmt.Printf("Latest PR comment is %d, updating.\n", commentID)
		return hub.UpdateComment(ctx, client, owner, repo, commentID, text)
	}

	fmt.Println("No PR comments, posting.")
	return hub.PostComment(ctx, client, owner, repo, number, text)
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
