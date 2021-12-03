package rpc

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
)

func (r *rpc) Webhooks(ctx *gin.Context, excludeRepos []string) {
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
		// terminate if repo substr URL is in the flags
		if excludeRepo(excludeRepos, *prEvent.Repo.URL) {
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

// Function that iterates over array of repo subString urls
func excludeRepo(excludeRepo []string, repoUrl string) bool {
	for _, repoUrlSubstr := range excludeRepo {

		return strings.Contains(repoUrl, repoUrlSubstr)

	}

	return false
}
