package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/skos-ninja/github-releaser/pkg/version"

	"github.com/google/go-github/v35/github"

	"github.com/leodido/go-conventionalcommits"
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
	log.Printf("PR: %v merged as %s with label %s\n", pr.GetNumber(), commitSHA, *versionType)

	latestVersion, err := getLatestVersion(ctx, client, repoOwner, repoName)
	if err != nil {
		return err
	}

	versionType, tagMessage = getVersionTypeAndTagMessage(ctx, client, repoOwner, repoName, pr.GetNumber(), latestVersion, pr.Title, pr.Labels, false)
	if versionType == nil {
		// PR does not have any of the following:
		// - a semantic title
		// - a valid label set
		return nil	
	}

	if tagMessage == nil {
		// TODO throw an error, the tag message was not created
		return nil
	}

	nextVerion, err := version.FindNextVersion(*versionType, latestVersion)
	if err != nil {
		return err
	}

	err = createTag(ctx, client, repoOwner, repoName, commitSHA, versionNum, tagMessage, pr.GetNumber(), a.impersonateTags)
	if err != nil {
		commentBody := fmt.Sprintf("Failed to make tag: `%s`", err.Error())
		createComment(ctx, client, pr.GetNumber(), repoOwner, repoName, commentBody)
	}
	return err
}

func getVersionTypeAndTagMessage(ctx context.Context, client *github.Client, repoOwner, repoName string, prNumber int, latestVersion string, prTitle string, prLabels []*github.Label, includePre bool) (*version.VersionType, string) {
	if versionType, parsedPrTitle := parseConventionalCommit(isDevelopmentVersion(latestVersion), prTitle); versionType != nil {
		var tagMessage := getTagMessageForConventionalCommitIncrement(parsedPrTitle)
		return versionType, tagMessage
	} else if versionType, _ := findVersionLabel(prLabels, includePre); versionType != nil{
		var tagMessage := getTagMessageForLabelBasedIncrement(ctx, client, repoOwner, repoName, pNumber)
		return versionType, tagMessage
	} else {
		return nil, nil
	}
}

func isDevelopmentVersion(version string) bool {
	return strings.HasPrefix(version, "v0")
}

func getTagMessageForConventionalCommitIncrement(parsedPrTitle ConventionalCommit) (string, error) {
	switch parsedPrTitle.Type {
	case "feat":
		return createTagMessage("features", parsedPrTitle.Description)
	case "fix":
		return createTagMessage("fixed", parsedPrTitle.Description)
    case "chore":
		return createTagMessage("chores", parsedPrTitle.Description)
	case "perf":
		return createTagMessage("improvements", parsedPrTitle.Description)
	default:
		// we map just the above ones to tag-police types: https://github.com/TrueLayer/tag-police/blob/master/tag_template.yml
		// if not mapped, just use the release notes and add the type
		return createTagMessage("release_notes", parsedPrTitle.Type + ": " + parsedPrTitle.Description)
    }
}

func createTagMessage(sectionName string, singleItem string) string {
	var message strings.Builder
	message.WriteString("%s:\n", sectionName)
	message.WriteString(fmt.Sprintf("  - \"%s\"\n", singleItem))
	return message.String()
}

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