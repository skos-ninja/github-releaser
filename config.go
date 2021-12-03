package main

import (
	. "github.com/skos-ninja/github-releaser/pkg/common"
)

var (
	cfg = &Config{
		Github:          Github{},
		ImpersonateTags: false,
		Port:            8080,
		ExcludeRepos:    []string{""},
	}
)
