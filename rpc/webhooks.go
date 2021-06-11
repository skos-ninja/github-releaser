package rpc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
)

func (r *rpc) Webhooks(ctx *gin.Context) {
	payload, err := github.ValidatePayload(ctx.Request, r.webhookSecretKey)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(ctx.Request), payload)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if prEvent, ok := event.(*github.PullRequestEvent); ok {
		err := r.app.HandleWebhook(ctx, prEvent)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.Status(http.StatusOK)
}
