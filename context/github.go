package context

import (
	"fmt"

	"github.com/Everlane/evan/config"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const ACCESS_TOKEN_FLAG string = "github.access_token"

// Memoized map of GitHub clients by their access token
var githubClients map[string]*github.Client

func (deployment *Deployment) GithubClient() (*github.Client, error) {
	if deployment.HasFlag(ACCESS_TOKEN_FLAG) {
		accessToken := deployment.Flag(ACCESS_TOKEN_FLAG).(string)
		if githubClients == nil {
			githubClients = make(map[string]*github.Client)
		}
		if githubClients[accessToken] == nil {
			githubClient := NewGithubClientWithAccessToken(accessToken)
			githubClients[accessToken] = githubClient
		}
		return githubClients[accessToken], nil
	}

	if config.DefaultGithubClient != nil {
		return config.DefaultGithubClient, nil
	}

	return nil, fmt.Errorf("No GitHub client configured nor was an access token found in '%v' flag", ACCESS_TOKEN_FLAG)
}

func NewGithubClientWithAccessToken(accessToken string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	return github.NewClient(tokenClient)
}
