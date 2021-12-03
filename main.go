package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
	"github.com/spf13/cobra"

	"github.com/skos-ninja/config-loader"
	"github.com/skos-ninja/github-releaser/app"
	. "github.com/skos-ninja/github-releaser/pkg/common"
	"github.com/skos-ninja/github-releaser/rpc"
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
	cmd.Flags().BoolP("impersonatetags", "i", cfg.ImpersonateTags, "Impersonate users when tagging")
	cmd.Flags().Int("appid", int(cfg.Github.AppId), "")
	cmd.Flags().String("webhooksecret", cfg.Github.WebhookSecretKey, "")
	cmd.Flags().String("privatekey", cfg.Github.PrivateKey, "")
	cmd.Flags().String("enterpriseurl", cfg.Github.EnterpriseURL, "")
	cmd.Flags().StringArray("excluderepos", cfg.ExcludeRepos, "List of repos to exclude")
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

	app, err := app.New(appTr, client, cfg.ImpersonateTags)
	if err != nil {
		return err
	}

	rpc := rpc.New(app, cfg.Github.WebhookSecretKey)

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Ok")
	})

	router.POST("/webhooks", func(c *gin.Context) {
		// get the '--excluderepos' flag value
		excludeRepos, _ := cmd.Flags().GetStringArray("excluderepos")

		rpc.Webhooks(c, excludeRepos)
	})

	return router.Run(fmt.Sprintf(":%v", cfg.Port))
}

func setupGithubClient(cfg Github) (*ghinstallation.AppsTransport, *github.Client, error) {
	tr := http.DefaultTransport

	// Due to envs vars loading \n as escaped we need to unescape
	privateKey := strings.ReplaceAll(cfg.PrivateKey, "\\n", "\n")
	itr, err := ghinstallation.NewAppsTransport(tr, cfg.AppId, []byte(privateKey))
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
