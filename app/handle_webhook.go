package app

import (
	"context"
	"log"

	"github.com/google/go-github/v35/github"
	"github.com/kr/pretty"
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
		err = a.handleClosed(ctx, client, prEvent)
	case "labeled":
		err = a.handleLabeled(ctx, client, prEvent)
	default:
		log.Printf("No handler for action: %s\n", action)
	}

	if err != nil {
		return handleError(err)
	}
	return nil
}

func handleError(err error) error {
	if e, ok := err.(*github.ErrorResponse); ok {
		pretty.Println(e.Errors)
		return e
	}

	return err
}
