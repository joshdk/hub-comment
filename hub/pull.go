// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package hub

import (
	"regexp"
	"strconv"
)

var (
	// rePullLink is a regex intended to match strings that look like
	// "https://github.com/joshdk/hub-comment/pull/123".
	rePullLink = regexp.MustCompile(`https://.*/([a-zA-A0-9_.-]+)/([a-zA-A0-9_.-]+)/pull/([1-9]\d*)`)
)

// splitPullRequestReference splits a pull request reference string into the
// owner, repo, and number of the PR.
//
// References should look like "https://github.com/joshdk/hub-comment/pull/123"
// or "joshdk/hub-comment#123".
func SplitPullRequestReference(reference string) (string, string, int, bool) {
	res := rePullLink.FindStringSubmatch(reference)
	if res == nil || len(res) != 4 {
		return "", "", 0, false
	}

	var (
		owner     = res[1]
		repo      = res[2]
		number, _ = strconv.Atoi(res[3])
	)

	return owner, repo, number, true
}
