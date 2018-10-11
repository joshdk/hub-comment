// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package hub

import (
	"context"

	"github.com/google/go-github/github"
)

// GetComments fetches all comments for the given pull request number.
func GetComments(ctx context.Context, client *github.Client, owner string, repo string, number int) ([]*github.IssueComment, error) {
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, number, nil)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// FilterComments selects the most recently updated comment, that was authored
// by the current user, if any exist.
func FilterComments(comments []*github.IssueComment, name string) (int64, bool) {
	var best *github.IssueComment
	for _, comment := range comments {
		user := comment.GetUser()
		if user == nil {
			continue
		}

		// Reject any comment that was not authored by the current user.
		author := user.GetLogin()
		if author != name {
			continue
		}

		// Keep any comment that has been more recently updated.
		if best == nil || comment.GetUpdatedAt().After(best.GetUpdatedAt()) {
			best = comment
		}
	}

	if best != nil {
		return *best.ID, true
	}
	return 0, false
}

// PostComment creates a new comment on the given PR.
func PostComment(ctx context.Context, client *github.Client, owner string, repo string, number int, comment string) (string, error) {
	ic := &github.IssueComment{
		Body: github.String(comment),
	}
	cmt, _, err := client.Issues.CreateComment(ctx, owner, repo, number, ic)
	if err != nil {
		return "", err
	}
	return cmt.GetHTMLURL(), nil
}

// UpdateComment modifies an existing comment.
func UpdateComment(ctx context.Context, client *github.Client, owner string, repo string, commentID int64, comment string) (string, error) {
	ic := &github.IssueComment{
		Body: github.String(comment),
	}
	cmt, _, err := client.Issues.EditComment(ctx, owner, repo, commentID, ic)
	if err != nil {
		return "", err
	}
	return cmt.GetHTMLURL(), nil
}
