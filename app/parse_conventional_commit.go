package app

import (
	"strings"

	conventionalcommitparser "github.com/release-lab/conventional-commit-parser"
	"github.com/skos-ninja/github-releaser/pkg/version"
)

func parseConventionalCommit(isDevelopmentVersion bool, text string) (*version.VersionType, conventionalcommitparser.Message) {
	parsedText := conventionalcommitparser.Parse(text)

	if parsedText.Header == "" {
		return nil, parsedText
	}

	var versionType string
	if isDevelopmentVersion {
		if strings.Contains(parsedText.Header, "!") {
			versionType = "MINOR"
		} else {
			versionType = "PATCH"
		}
	} else {
		if strings.Contains(parsedText.Header, "!") {
			versionType = "MAJOR"
		} else if parsedText.ParseHeader().Scope == "feat" {
			versionType = "MINOR"
		} else {
			versionType = "PATCH"
		}
	}

	parsedVersionType, error := version.ParseVersionType(versionType)

	if error != nil {
		return nil, parsedText
	}

	return &parsedVersionType, parsedText
}
