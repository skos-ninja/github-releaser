package app

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v35/github"
)

func getLatestVersion(ctx context.Context, client *github.Client, repoOwner, repoName string) (string, error) {
	refs, _, err := client.Git.ListMatchingRefs(ctx, repoOwner, repoName, &github.ReferenceListOptions{
		Ref: "tags",
	})
	if err != nil {
		return "", err
	}

	versions := []*semver.Version{}
	for _, ref := range refs {
		tag := strings.TrimPrefix(ref.GetRef(), "refs/tags/")
		v, err := semver.NewVersion(tag)
		if err != nil {
			log.Printf("invalid semver %s: %s\n", tag, err.Error())
			continue
		}

		t := ref.GetObject().GetType()
		if t != "commit" && t != "tag" {
			log.Printf("invalid ref %s: %s\n", tag, t)
			continue
		}

		versions = append(versions, v)
	}

	sort.Sort(semver.Collection(versions))
	if len(versions) == 0 {
		return "v0.0.0", nil
	}

	version := versions[len(versions)-1].Original()
	return version, nil
}
