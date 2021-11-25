package app

import (
	"context"
	"fmt"
	"strings"


	"github.com/google/go-github/v35/github"
)

func getTagMessageForLabelBasedIncrement(ctx context.Context, client *github.Client, repoOwner, repoName string, prNumber int) (string, error) {
	var message strings.Builder
	commits, _, err := client.PullRequests.ListCommits(ctx, repoOwner, repoName, prNumber, &github.ListOptions{})
	if err != nil {
		return "", err
	}
	message.WriteString("release_notes:\n")
	for _, commit := range commits {
		if cm := commit.Commit.GetMessage(); cm != "" {
			message.WriteString(fmt.Sprintf("  - [%s] %s\n", commit.GetSHA()[:9], cm))
		}
	}

	return message.String(), nil
}