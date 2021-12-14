package rpc

import (
	"github.com/gin-gonic/gin"

	"github.com/skos-ninja/github-releaser/app"
)

type RPC interface {
	Webhooks(ctx *gin.Context)
}

type rpc struct {
	app              app.App
	webhookSecretKey []byte
	excludedRepos    []string
}

func New(app app.App, webhookSecretKey string, excludedRepos []string) RPC {
	return &rpc{app, []byte(webhookSecretKey), excludedRepos}
}
