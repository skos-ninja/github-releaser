package rpc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v41/github"
	"github.com/skos-ninja/github-releaser/pkg/common"
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
		// terminate if repo name is set to be excluded
		if common.Contains(r.excludedRepos, prEvent.Repo.GetFullName()) {
			return
		}

		err := r.app.HandleWebhook(ctx, prEvent)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.Status(http.StatusOK)
}
