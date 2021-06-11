package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/skos-ninja/github-releaser/pkg/version"

	"github.com/google/go-github/v35/github"
	"github.com/kr/pretty"
)

func handleClosed(ctx context.Context, client *github.Client, prEvent *github.PullRequestEvent) error {
	pr := prEvent.GetPullRequest()
	if pr == nil {
		// Ignoring as we are missing pull request data.
		return errors.New("Missing pull request data")
	}

	if !pr.GetMerged() {
		// Ignoring as the pr was not merged
		return nil
	}

	versionType := findVersionLabel(pr.Labels)
	if versionType == nil {
		// PR does not have a valid label set
		return nil
	}

	commitSHA := pr.GetMergeCommitSHA()
	log.Printf("PR: %v merged as %s with label %s\n", pr.GetNumber(), commitSHA, *versionType)

	repo := prEvent.GetRepo()
	if repo == nil {
		return errors.New("Missing repository information")
	}
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()

	ver, err := getLatestVersion(ctx, client, repoOwner, repoName)
	if err != nil {
		return handleError(err)
	}

	versionNum, err := version.FindNextVersion(*versionType, ver)
	if err != nil {
		return err
	}

	tagURL, err := createTag(ctx, client, repoOwner, repoName, commitSHA, versionNum)
	if err != nil {
		return handleError(err)
	}

	commentBody := fmt.Sprintf("[%s](%s) created on %s", versionNum, tagURL, commitSHA)
	comment := &github.IssueComment{
		Body: &commentBody,
	}
	_, _, err = client.Issues.CreateComment(ctx, repoOwner, repoName, pr.GetNumber(), comment)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func findVersionLabel(labels []*github.Label) *version.VersionType {
	p := 0
	var v *version.VersionType = nil
	for _, label := range labels {
		nV, err := version.ParseVersionType(strings.ToUpper(label.GetName()))
		if err != nil || nV == version.Prerelease {
			// We shouldn't release a tag on a merged PR with the label pre-release
			continue
		}

		nP := version.GetVersionPriority(nV)
		if nP > p {
			p = nP
			v = &nV
		}
	}

	return v
}

func handleError(err error) error {
	if e, ok := err.(*github.ErrorResponse); ok {
		pretty.Println(e.Errors)
		return e
	}

	return err
}
