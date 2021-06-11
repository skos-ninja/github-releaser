package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/skos-ninja/github-releaser/app"
	"github.com/skos-ninja/github-releaser/rpc"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
	"github.com/skos-ninja/config-loader"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:  "github-releaser",
		Args: cobra.ExactArgs(0),
		RunE: runE,
	}
)

func init() {
	cmd.AddCommand(bumpVersionCmd)

	config.Init(cmd)
	cmd.Flags().IntP("port", "p", cfg.Port, "HTTP Listening port")
	cmd.Flags().Int("appid", int(cfg.Github.AppId), "")
	cmd.Flags().String("webhooksecret", cfg.Github.WebhookSecretKey, "")
	cmd.Flags().String("privatekey", cfg.Github.PrivateKey, "")
	cmd.Flags().String("enterpriseurl", cfg.Github.EnterpriseURL, "")
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func runE(cmd *cobra.Command, args []string) error {
	if err := config.Load(cmd, cfg); err != nil {
		return err
	}

	appTr, client, err := setupGithubClient(cfg.Github)
	if err != nil {
		return err
	}
	app, err := app.New(appTr, client)
	if err != nil {
		return err
	}
	rpc := rpc.New(app, cfg.Github.WebhookSecretKey)

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Ok")
	})

	router.POST("/webhooks", rpc.Webhooks)

	return router.Run(fmt.Sprintf(":%v", cfg.Port))
}

func setupGithubClient(cfg Github) (*ghinstallation.AppsTransport, *github.Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewAppsTransport(tr, cfg.AppId, []byte(cfg.PrivateKey))
	if err != nil {
		return nil, nil, err
	}

	var client *github.Client
	if cfg.EnterpriseURL != "" {
		fmt.Printf("Using enterprise: %s\n", cfg.EnterpriseURL)
		itr.BaseURL = cfg.EnterpriseURL
		client, err = github.NewEnterpriseClient(cfg.EnterpriseURL, cfg.EnterpriseURL, &http.Client{Transport: itr})
		if err != nil {
			return itr, nil, err
		}
	} else {
		client = github.NewClient(&http.Client{Transport: itr})
	}

	return itr, client, nil
}
