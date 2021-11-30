package rpc

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
)

var (
	excludeRepos = []string{
		"https://github.com/TrueLayer/prisma-tbd",
	}
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
		if contains(excludeRepos, *prEvent.Repo.URL) {
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
