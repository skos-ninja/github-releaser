package app

import (
	"strings"

	"github.com/skos-ninja/github-releaser/pkg/version"

	"github.com/google/go-github/v35/github"
)

func findVersionLabel(labels []*github.Label, includePre bool) (*version.VersionType, string) {
	p := 0
	var v *version.VersionType = nil
	var l = ""
	for _, label := range labels {
		nV, err := version.ParseVersionType(strings.ToUpper(label.GetName()))
		if err != nil || (!includePre && nV == version.Prerelease) {
			// We shouldn't release a tag on a merged PR with the label pre-release
			continue
		}

		nP := version.GetVersionPriority(nV)
		if nP > p {
			p = nP
			v = &nV
			l = label.GetName()
		}
	}

	return v, l
}
