package common

import (
	"github.com/google/go-github/github"
)

type DeploymentState int

const (
	DEPLOYMENT_PENDING DeploymentState = iota
	RUNNING_PRECONDITIONS
	RUNNING_PHASE
	DEPLOYMENT_DONE
	DEPLOYMENT_ERROR
)

type DeploymentStatus struct {
	State DeploymentState
	Phase Phase
	Error error
}

type Deployment interface {
	Application() Application
	Environment() string
	Ref() string
	GithubClient() *github.Client
	Flags() map[string]interface{}
	IsForce() bool

	Status() DeploymentStatus
}

type Application interface {
	Name() string
	Repository() Repository
}

type Repository interface {
	Owner() string
	Name() string
}

func CanonicalNameForRepository(repository Repository) string {
	owner := repository.Owner()
	name := repository.Name()
	return owner + "/" + name
}
