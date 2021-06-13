package app

import (
	"context"
	"errors"
	"log"

	"github.com/skos-ninja/github-releaser/pkg/version"

	"github.com/google/go-github/v35/github"
)

func handleLabeled(ctx context.Context, client *github.Client, prEvent *github.PullRequestEvent) error {
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

	versionType, labelName := findVersionLabel(pr.Labels, true)
	if versionType == nil {
		// PR does not have a valid label set
		log.Println("PR has no valid label set")
		return nil
	} else if *versionType != version.Prerelease {
		// PR doesn't have a valid label set for this event
		log.Printf("PR label not of prerelease. Is %s\n", *versionType)
		return nil
	}

	head := pr.GetHead()
	if head == nil {
		return errors.New("missing pr head data")
	}
	headRepo := head.GetRepo()
	if headRepo == nil {
		return errors.New("missing pr head repo")
	}
	headOwner := headRepo.GetOwner().GetLogin()
	headBranch := head.GetRef()

	if headOwner != repoOwner {
		commentBody := "PR branch needs to be created on original repo for tagging"
		err := createComment(ctx, client, pr.GetNumber(), repoOwner, repoName, commentBody)
		if err != nil {
			return err
		}
	}

	v, err := getLatestBranchVersion(ctx, client, repoOwner, repoName, headBranch)
	if err != nil {
		return err
	}
	versionNum, err := version.FindNextVersion(*versionType, v)
	if err != nil {
		return err
	}

	commitSHA := head.GetSHA()
	prNum := pr.GetNumber()
	err = createTag(ctx, client, repoOwner, repoName, commitSHA, versionNum, repo.GetHTMLURL(), prNum)
	if err != nil {
		return err
	}

	_, err = client.Issues.RemoveLabelForIssue(ctx, repoOwner, repoName, prNum, labelName)
	return err
}
