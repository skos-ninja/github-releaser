package main

import (
	"log"
	"strings"

	"github.com/skos-ninja/github-releaser/pkg/version"

	"github.com/spf13/cobra"
)

var (
	bumpVersionCmd = &cobra.Command{
		Use:   "bump-version <semver> <major|minor|patch|prerelease>",
		Short: "Command line version bumping",
		Long:  "Input a semver version with the version bump type you want and the version will be output in stdout",
		Args:  cobra.ExactArgs(2),
		RunE:  bumpVersionRunE,
	}
)

func bumpVersionRunE(cmd *cobra.Command, args []string) error {
	v := args[0]
	t, err := version.ParseVersionType(strings.ToUpper(args[1]))
	if err != nil {
		return err
	}
	log.Printf("Bumping version %s by %s\n", v, t)

	ver, err := version.FindNextVersion(t, v)
	if err != nil {
		return err
	}

	log.Printf("Version: %s", ver)
	return nil
}
