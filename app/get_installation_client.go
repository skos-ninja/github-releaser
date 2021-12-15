package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	ghinstallation "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
)

type installLruKey string

func (a *app) GetInstallationClient(ctx context.Context, org, user string) (*github.Client, error) {
	if v, ok := a.installLru.Get(installLruKey(org + user)); ok {
		log.Println("Using cached client")
		if client, ok := v.(*github.Client); ok {
			return client, nil
		}
	}

	var install *github.Installation = nil
	var err error
	if org != "" {
		install, _, err = a.github.Apps.FindOrganizationInstallation(ctx, org)
	} else if user != "" {
		install, _, err = a.github.Apps.FindUserInstallation(ctx, user)
	} else {
		err = errors.New("user or org not provided")
	}
	if err != nil {
		return nil, err
	}

	id := install.GetID()
	if id == 0 {
		return nil, errors.New("app not installed")
	}

	itr := ghinstallation.NewFromAppsTransport(a.appTr, id)
	client := github.NewClient(&http.Client{Transport: itr})
	a.installLru.Add(installLruKey(org+user), client)
	log.Printf("Stored client: %s", org+user)

	return client, nil
}
