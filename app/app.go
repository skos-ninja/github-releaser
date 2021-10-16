package app

import (
	"context"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v35/github"
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
