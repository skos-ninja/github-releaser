package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v35/github"
)

func createTag(ctx context.Context, client *github.Client, repoOwner, repoName, commitSHA, version string, prNumber int, impersonateTags bool) error {
	message, err := getTagMessage(ctx, client, repoOwner, repoName, prNumber)
	if err != nil {
		return handleError(err)
	}

	var tagger *github.CommitAuthor
	if impersonateTags {
		commit, _, err := client.Git.GetCommit(ctx, repoOwner, repoName, commitSHA)
		if err != nil {
			return handleError(err)
		}
		tagger = commit.Author
	}

	tag := &github.Tag{
		Tag:     &version,
		SHA:     &commitSHA,
		Message: github.String(message),
		Tagger:  tagger,
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  &commitSHA,
		},
	}
	t, _, err := client.Git.CreateTag(ctx, repoOwner, repoName, tag)
	if err != nil {
		return handleError(err)
	}

	ref := &github.Reference{
		Ref: github.String(fmt.Sprintf("refs/tags/%s", version)),
		Object: &github.GitObject{
			SHA: t.SHA,
		},
	}
	_, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, ref)
	if err != nil {
		log.Printf("Failed to make ref: %s\n", version)
		return handleError(err)
	}

	commentBody := fmt.Sprintf("[%s](../releases/tag/%s) created on %s", version, version, commitSHA)
	err = createComment(ctx, client, prNumber, repoOwner, repoName, commentBody)
	return err
}

func getTagMessage(ctx context.Context, client *github.Client, repoOwner, repoName string, prNumber int) (string, error) {
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
