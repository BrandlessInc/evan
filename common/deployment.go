package common

import (
	"github.com/google/go-github/github"
)

type Deployment interface {
	Application() Application
	Ref() string
	GithubClient() *github.Client
	Flags() map[string]interface{}
	IsForce() bool

	// Some object representing the request that initiated this deployment,
	// eg. `*github.DeploymentEvent`.
	Initiator() interface{}
}

type Application interface {
	Repository() Repository
}

type Repository interface {
	Owner() string
	Name() string
}
