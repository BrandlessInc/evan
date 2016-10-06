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

func (deployment *Deployment) GetApplication() *config.Application {
	return deployment.Application
}

func (deployment *Deployment) GetGithubClient() *github.Client {
	return deployment.GithubClient
}

func (deployment *Deployment) GetRef() string {
	return deployment.Ref
}

func (deployment *Deployment) RunPreconditions() []config.PreconditionResult {
	preconditions := deployment.Strategy.Preconditions

	resultChan := make(chan config.PreconditionResult)
	for _, precondition := range preconditions {
		go precondition.Status(deployment, resultChan)
	}

	results := make([]config.PreconditionResult, 0, len(preconditions))
	for i := range preconditions {
		results[i] = <-resultChan
	}
	return results
}

func (deployment *Deployment) Run() error {
	for _, result := range deployment.RunPreconditions() {
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
