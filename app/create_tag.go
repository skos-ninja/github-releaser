package app

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v35/github"
)

func createTag(ctx context.Context, client *github.Client, repoOwner, repoName, commit, version, repoURL string, prNumber int) error {
	tag := &github.Tag{
		Tag:     &version,
		SHA:     &commit,
		Message: github.String(""),
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  &commit,
		},
	}
	t, _, err := client.Git.CreateTag(ctx, repoOwner, repoName, tag)
	if err != nil {
		return handleError(err)
	}

	ref := &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/tags/%s", version)),
		Object: t.Object,
	}
	_, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, ref)
	if err != nil {
		log.Printf("Failed to make ref: %s\n", version)
		return handleError(err)
	}

	tagURL := fmt.Sprintf("%s/releases/tag/%s", repoURL, version)
	commentBody := fmt.Sprintf("[%s](%s) created on %s", version, tagURL, commit)
	err = createComment(ctx, client, prNumber, repoOwner, repoName, commentBody)
	return err
}
