package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/skos-ninja/github-releaser/pkg/version"

	conventionalcommitparser "github.com/release-lab/conventional-commit-parser"

	"github.com/google/go-github/v35/github"
)

func (a *app) handleClosed(ctx context.Context, client *github.Client, prEvent *github.PullRequestEvent) error {
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

	commitSHA := pr.GetMergeCommitSHA()
	log.Printf("PR: %v merged as %s\n", pr.GetNumber(), commitSHA)

	latestVersion, err := getLatestVersion(ctx, client, repoOwner, repoName)
	if err != nil {
		return err
	}

	versionType, tagMessage, error := getVersionTypeAndTagMessage(ctx, client, repoOwner, repoName, pr.GetNumber(), latestVersion, *pr.Title, pr.Labels, false)
	if error != nil {
		return nil
	}
	if versionType == nil {
		// PR does not have any of the following:
		// - a semantic title
		// - a valid label set
		commentBody := fmt.Sprintf("Github-releaser is installed in this repository but could not create any tag. PR title is not semantic and there are no labels.")
		createComment(ctx, client, pr.GetNumber(), repoOwner, repoName, commentBody)
		return nil
	}

	if tagMessage == "" {
		// TODO throw an error, the tag message was not created
		return nil
	}

	nextVersion, err := version.FindNextVersion(*versionType, latestVersion)
	if err != nil {
		return err
	}

	err = createTag(ctx, client, repoOwner, repoName, commitSHA, nextVersion, tagMessage, pr.GetNumber(), a.impersonateTags)
	if err != nil {
		commentBody := fmt.Sprintf("Failed to make tag: `%s`", err.Error())
		createComment(ctx, client, pr.GetNumber(), repoOwner, repoName, commentBody)
	}
	return err
}

func getVersionTypeAndTagMessage(ctx context.Context, client *github.Client, repoOwner, repoName string, prNumber int, latestVersion string, prTitle string, prLabels []*github.Label, includePre bool) (*version.VersionType, string, error) {
	if versionType, parsedPrTitle := parseConventionalCommit(isDevelopmentVersion(latestVersion), prTitle); versionType != nil {
		tagMessage := getTagMessageForConventionalCommitIncrement(parsedPrTitle, prNumber)
		return versionType, tagMessage, nil
	} else if versionType, _ := findVersionLabel(prLabels, includePre); versionType != nil {
		tagMessage, error := getTagMessageForLabelBasedIncrement(ctx, client, repoOwner, repoName, prNumber)
		return versionType, tagMessage, error
	} else {
		return nil, "", nil
	}
}

func isDevelopmentVersion(version string) bool {
	return strings.HasPrefix(version, "v0")
}

func getTagMessageForConventionalCommitIncrement(parsedPrTitle conventionalcommitparser.Message, prNumber int) string {
	commitType := parsedPrTitle.ParseHeader().Type
	commitScope := parsedPrTitle.ParseHeader().Scope
	commitDescription := parsedPrTitle.ParseHeader().Subject

	switch commitType {
	case "feat":
		return createTagMessage("features", "("+commitScope+"): "+commitDescription, prNumber)
	case "fix":
		return createTagMessage("fixed", "("+commitScope+"): "+commitDescription, prNumber)
	case "chore":
		return createTagMessage("chores", "("+commitScope+"): "+commitDescription, prNumber)
	case "perf":
		return createTagMessage("improvements", "("+commitScope+"): "+commitDescription, prNumber)
	default:
		// we map just the above ones to tag-police types: https://github.com/TrueLayer/tag-police/blob/master/tag_template.yml
		// if not mapped, just use the release notes and add the type
		return createTagMessage("release_notes", commitType+"("+commitScope+"): "+commitDescription, prNumber)
	}
}

func createTagMessage(sectionName string, singleItem string, prNumber int) string {
	var message strings.Builder
	message.WriteString(fmt.Sprintf("%s:\n", sectionName))
	message.WriteString(fmt.Sprintf("  - \"%s [#%s]\"\n", singleItem, prNumber))
	return message.String()
}
