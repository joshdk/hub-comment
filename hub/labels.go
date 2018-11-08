// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package hub

import (
	"context"
	"sort"

	"github.com/google/go-github/github"
)

// GetLabels fetches all labels for the given pull request number.
func GetLabels(ctx context.Context, client *github.Client, owner string, repo string, number int) ([]*github.Label, error) {
	labels, _, err := client.Issues.ListLabelsByIssue(ctx, owner, repo, number, nil)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// OnlyLabelNames simplifies a list of GitHub labels into a list a strings,
// which is then sorted alphabetically.
func OnlyLabelNames(labels []*github.Label) []string {
	list := make([]string, len(labels))
	for index, label := range labels {
		list[index] = label.GetName()
	}
	sort.Strings(list)
	return list
}
