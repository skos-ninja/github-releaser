package app

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func createComment(ctx context.Context, client *github.Client, prNum int, repoOwner, repoName, comment string) error {
	issueComment := &github.IssueComment{
		Body: &comment,
	}
	_, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, issueComment)
	return err
}
