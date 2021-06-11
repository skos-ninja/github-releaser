package rpc

import (
	"github.com/skos-ninja/github-releaser/app"

	"github.com/gin-gonic/gin"
)

type RPC interface {
	Webhooks(ctx *gin.Context)
}

type rpc struct {
	app              app.App
	webhookSecretKey []byte
}

func New(app app.App, webhookSecretKey string) RPC {
	return &rpc{app, []byte(webhookSecretKey)}
}
