package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/skos-ninja/github-releaser/pkg/version"

	"github.com/google/go-github/v35/github"
)

func handleClosed(ctx context.Context, client *github.Client, prEvent *github.PullRequestEvent) error {
	pr := prEvent.GetPullRequest()
	if pr == nil {
		// Ignoring as we are missing pull request data.
		return errors.New("missing pull request data")
	}

	repo := prEvent.GetRepo()
	if repo == nil {
		return errors.New("missing repository information")
	}
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()

	if !pr.GetMerged() || repo.GetDefaultBranch() != pr.GetBase().GetRef() {
		// Ignoring as the pr was not merged into default
		return nil
	}

	versionType, _ := findVersionLabel(pr.Labels, false)
	if versionType == nil {
		// PR does not have a valid label set
		return nil
	}

	commitSHA := pr.GetMergeCommitSHA()
	log.Printf("PR: %v merged as %s with label %s\n", pr.GetNumber(), commitSHA, *versionType)

	ver, err := getLatestVersion(ctx, client, repoOwner, repoName)
	if err != nil {
		return err
	}

	versionNum, err := version.FindNextVersion(*versionType, ver)
	if err != nil {
		return err
	}

	err = createTag(ctx, client, repoOwner, repoName, commitSHA, versionNum, pr.GetNumber())
	if err != nil {
		commentBody := fmt.Sprintf("Failed to make tag: `%s`", err.Error())
		createComment(ctx, client, pr.GetNumber(), repoOwner, repoName, commentBody)
	}
	return err
}
