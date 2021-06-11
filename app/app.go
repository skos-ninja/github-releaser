package app

import (
	"context"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v35/github"
	"github.com/hashicorp/golang-lru/simplelru"
)

type App interface {
	GetZen(ctx context.Context) (string, error)

	HandleWebhook(ctx context.Context, pr *github.PullRequestEvent) error
}

type app struct {
	appTr      *ghinstallation.AppsTransport
	github     *github.Client
	installLru *simplelru.LRU
}

func New(appTr *ghinstallation.AppsTransport, github *github.Client) (App, error) {
	iLRU, err := simplelru.NewLRU(100, nil)
	if err != nil {
		return nil, err
	}

	return &app{
		appTr,
		github,
		iLRU,
	}, nil
}
