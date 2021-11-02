package app

import (
	"strings"
	
	"github.com/leodido/go-conventionalcommits"
)

func parseConventionalCommit(isDevelopmentVersion bool, text string) (*version.VersionType, ConventionalCommit) {
	parsedText, error := parser.NewMachine(WithTypes(conventionalcommits.TypesConventional)).Parse(text)

	if error != nil {
		// there was an error when parsing, the title does not follow the spec
		return nil
	}
	
	var versionType
	if isDevelopmentVersion {
		if parsedText.IsBreakingChange() {
			versionType := "MINOR"
		} else {
			versionType := "PATCH"
		}
	} else {
		if parsedText.IsBreakingChange() {
			versionType := "MAJOR"
		} else if parsedText.Type == "feat" {
			versionType := "MINOR"
		} else {
			versionType := "PATCH"
		}
	}

	return versionType, parsedText
}
