package rpc

import (
	"github.com/gin-gonic/gin"

	"github.com/skos-ninja/github-releaser/app"
)

type RPC interface {
	Webhooks(ctx *gin.Context, excludeRepos []string)
}

type rpc struct {
	app              app.App
	webhookSecretKey []byte
}

func New(app app.App, webhookSecretKey string) RPC {
	return &rpc{app, []byte(webhookSecretKey)}
}
