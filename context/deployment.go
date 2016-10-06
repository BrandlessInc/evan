package context

import (
	"github.com/Everlane/evan/config"

	"github.com/google/go-github/github"
)

// Stores state relating to a deployment.
type Deployment struct {
	Application  *config.Application
	Environment  string
	Strategy     *config.Strategy
	Ref          string
	GithubClient *github.Client
}

func (deployment *Deployment) GetGithubClient() *github.Client {
	return deployment.GithubClient
}

func (deployment *Deployment) GetRef() string {
	return deployment.Ref
}
