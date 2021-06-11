package app

import (
	"context"
	"fmt"

	"github.com/google/go-github/v35/github"
)

func createTag(ctx context.Context, client *github.Client, repoOwner, repoName, commit, version string) error {
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
		return handleError(err)
	}

	return nil
}
