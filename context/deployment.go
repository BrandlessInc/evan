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
