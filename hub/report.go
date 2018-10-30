package hub

import (
	"fmt"
	"strings"
)

// Report display a textual report about the pull request, and comment that was
// just posted.
func Report(comment, url, user, login, owner, repo string, number int, updated bool, dryRun bool) {
	var (
		prefix = "Posting new comment as"
		lines  = strings.Split(comment, "\n")
	)

	switch {
	case updated && dryRun:
		prefix = "Would have updated existing comment by"
	case updated:
		prefix = "Updating existing comment by"
	case dryRun:
		prefix = "Would have posted new comment as"
	}

	fmt.Printf(
		"%s %s (%s) on %s/%s#%d:\n",
		prefix,
		user, login,
		owner, repo, number,
	)

	fmt.Println()
	for _, line := range lines {
		fmt.Printf("→ %s\n", line)
	}

	if dryRun {
		return
	}

	fmt.Println()
	fmt.Println("To view comment visit:")
	fmt.Printf("→ %s\n", url)
}
