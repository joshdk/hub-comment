package hub

import (
	"fmt"
	"strings"
)

// Report display a textual report about the pull request, and comment that was
// just posted.
func Report(comment, url, user, login, owner, repo string, number int, updated bool) {
	var (
		prefix = "Posting new comment as"
		lines  = strings.Split(comment, "\n")
	)

	if updated {
		prefix = "Updating existing comment by"
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
	fmt.Println()

	fmt.Println("To view comment visit:")
	fmt.Printf("→ %s\n", url)
}
