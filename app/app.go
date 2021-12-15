package app

import (
	"context"

	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
	"github.com/hashicorp/golang-lru/simplelru"
)

type App interface {
	HandleWebhook(ctx context.Context, pr *github.PullRequestEvent) error
}

type app struct {
	appTr           *ghinstallation.AppsTransport
	github          *github.Client
	installLru      *simplelru.LRU
	impersonateTags bool
}

func New(appTr *ghinstallation.AppsTransport, github *github.Client, impersonateTags bool) (App, error) {
	iLRU, err := simplelru.NewLRU(100, nil)
	if err != nil {
		return nil, err
	}

	return &app{
		appTr,
		github,
		iLRU,
		impersonateTags,
	}, nil
}
