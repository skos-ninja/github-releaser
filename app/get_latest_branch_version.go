package app

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v41/github"
	"golang.org/x/sync/errgroup"
)

func getLatestBranchVersion(ctx context.Context, client *github.Client, repoOwner, repoName, branch string) (string, error) {
	var refs []*github.Reference
	var commits []*github.RepositoryCommit
	g, gctx := errgroup.WithContext(ctx)

	// Tags
	g.Go(func() (err error) {
		refs, _, err = client.Git.ListMatchingRefs(gctx, repoOwner, repoName, &github.ReferenceListOptions{
			Ref: "tags",
		})

		return
	})

	// Commits
	g.Go(func() (err error) {
		commits, _, err = client.Repositories.ListCommits(gctx, repoOwner, repoName, &github.CommitsListOptions{
			SHA: branch,
		})

		return
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	versions := []*semver.Version{}
	repoVersions := []*semver.Version{}
	for _, ref := range refs {
		tag := strings.TrimPrefix(ref.GetRef(), "refs/tags/")
		v, err := semver.NewVersion(tag)
		if err != nil {
			log.Printf("invalid semver %s: %s\n", tag, err.Error())
			continue
		}

		t := ref.GetObject().GetType()
		if t == "commit" {
			sha := ref.GetObject().GetSHA()
			for _, c := range commits {
				if sha == c.GetSHA() {
					repoVersions = append(repoVersions, v)
					break
				}
			}
		} else if v.Prerelease() == "" && t == "tag" {
			repoVersions = append(repoVersions, v)
		} else {
			log.Printf("invalid ref %s: %s\n", tag, t)
			continue
		}

		versions = append(versions, v)
	}

	var vs semver.Collection = versions
	if len(repoVersions) != 0 {
		log.Println("Using tags on branch only")
		vs = repoVersions
	}

	sort.Sort(vs)
	if len(vs) == 0 {
		return "v0.0.0", nil
	}

	version := vs[len(vs)-1].Original()
	return version, nil
}
