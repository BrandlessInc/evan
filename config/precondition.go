package config

import (
	"github.com/google/go-github/github"
)

type PreconditionResult struct {
	Precondition Precondition
	Error        error
}

func createResult(precondition Precondition, err error) PreconditionResult {
	return PreconditionResult{
		Precondition: precondition,
		Error:        err,
	}
}

// For parallelization preconditions report their statuses on a channel.
type PreconditionResults chan PreconditionResult

type Precondition interface {
	Status(*Application, Deployment, PreconditionResults)
}

type Deployment interface {
	GetRef() string
	GetGithubClient() *github.Client
}
