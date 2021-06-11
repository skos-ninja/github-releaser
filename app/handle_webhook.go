package app

import (
	"context"
	"log"

	"github.com/google/go-github/v35/github"
)

func (a *app) HandleWebhook(ctx context.Context, prEvent *github.PullRequestEvent) error {
	prOrg := ""
	prUser := ""
	if org := prEvent.GetOrganization(); org != nil {
		prOrg = org.GetLogin()
	} else {
		prUser = prEvent.GetRepo().GetOwner().GetLogin()
	}
	client, err := a.GetInstallationClient(ctx, prOrg, prUser)
	if err != nil {
		return err
	}

	action := prEvent.GetAction()
	switch action {
	case "closed":
		return handleClosed(ctx, client, prEvent)
	}

	log.Printf("No handler for action: %s\n", action)
	return nil
}
