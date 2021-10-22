# github-releaser

This tool is designed to perform automated semver tagging of pull requests.

## Instalation

The install is done in two parts. You need to make the docker image accessible on a url that is public to at least GitHubs IPs

### Docker

The docker image expects you to provide at least:
- `GITHUB_APP_ID`
- `GITHUB_WEBHOOK_SECRET`
- `GITHUB_PRIVATE_KEY`
and optionally `GITHUB_ENTERPRISE_URL`

Once running the container will expose the endpoint on port `8080` (or can be set by `HTTP_PORT`) and the webhook receiver url is

### Github
To install this app you will need to create a GitHub app and assign it a webhook url for 

#### Permissions
- Contents (Read & Write)
- Metadata (Read-only)
- Pull Requests (Read & Write)

#### Events
- Pull Request