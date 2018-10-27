// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
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

// version can be replaced at build time with a custom version string.
var version = "development"

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "hub-comment: %s\n", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	var (
		// versionFlag is a command line flag ("-version") that will have the
		// program to print a version string and then exit.
		versionFlag = flag.Bool("version", false, fmt.Sprintf(`Print the version "%s" and exit.`, version))

		// templateFlag is a command line flag ("-template") that holds a
		// string literal used as the posted comment body.
		templateFlag = flag.String("template", "", "Comment body to post.")

		// templateFileFlag is a command line flag ("-template-file") that
		// names a file, the contents of which is used as the posted comment
		// body.
		templateFileFlag = flag.String("template-file", "", "File containing comment body to post.")

		// typeFlag is a command line flag ("-type") that specifies the type of
		// comment to post. Comment types are completely arbitrary, but are
		// used from distinguishing between multiple different kinds of
		// comments on a single PR.
		typeFlag = flag.String("type", "default", "Type of comment to post and edit.")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of hub-comment:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return nil
	}

	token, commentExists := os.LookupEnv(githubTokenEnvVar)
	if !commentExists {
		return fmt.Errorf("no GITHUB_TOKEN set in environment")
	}

	// If no CIRCLE_PULL_REQUEST is set, print an error and return immediately
	// but do not fail. This environment variable will not be set on non-pr
	// branches, or if a build is started before a pr is opened.
	reference, commentExists := os.LookupEnv(pullRequestLinkEnvVar)
	if !commentExists {
		fmt.Fprintln(os.Stderr, "hub-comment: no CIRCLE_PULL_REQUEST set in environment")
		return nil
	}

	owner, repo, number, commentExists := hub.SplitPullRequestReference(reference)
	if !commentExists {
		return fmt.Errorf("malformed pull request link")
	}

	// Get a template from either the -template flag directly, or read from the
	// -template-file.
	template, err := getTemplate(*templateFlag, *templateFileFlag)
	if err != nil {
		return err
	}

	// Parse the template body
	tpl, err := hub.NewTemplate(template)
	if err != nil {
		return err
	}

	// Build a context object containing the available environment variables.
	state := hub.NewContext(os.Environ(), *typeFlag)

	comment, err := hub.Execute(tpl, state)
	if err != nil {
		return err
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

	// Get a list of all comments for the given PR number
	comments, err := hub.GetComments(ctx, client, owner, repo, number)
	if err != nil {
		return err
	}

	// Select the most recent comment that was authored by the current user, if
	// one exists.
	commentID, commentExists := hub.FilterComments(comments, self.GetLogin(), *typeFlag)

	// Create a new comment or update an existing comment. Save a link to the
	// resulting comment.
	var url string
	if commentExists {
		url, err = hub.UpdateComment(ctx, client, owner, repo, commentID, comment)
	} else {
		url, err = hub.PostComment(ctx, client, owner, repo, number, comment)
	}
	if err != nil {
		return err
	}

	// Display a report about the comment that was just posted.
	hub.Report(comment, url, self.GetName(), self.GetLogin(), owner, repo, number, commentExists)

	return nil
}

// getTemplate will either return the contents of template verbatim, or return
// the contents read from templateFile.
func getTemplate(template string, templateFile string) ([]byte, error) {
	switch {
	case template == "" && templateFile == "":
		fallthrough
	case template != "" && templateFile != "":
		return nil, fmt.Errorf("a template or a template file must be given")
	case template != "":
		return []byte(template), nil
	default:
		return ioutil.ReadFile(templateFile)
	}
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
