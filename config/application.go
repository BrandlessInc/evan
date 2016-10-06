package config

import (
	"github.com/google/go-github/github"
)

type Application struct {
	Repository   *Repository
	Environments []string

	// Called to determine the target and strategy to use for deploying to a
	// given environment.
	DeployEnvironment func(string) *Strategy
}

type Applications struct {
	Map map[string]*Application
}

func (apps *Applications) FindApplicationForGithubRepository(repo *github.Repository) *Application {
	for _, app := range apps.Map {
		if app.Repository.Owner == *repo.Owner.Login && app.Repository.Name == *repo.Name {
			return app
		}
	}

	return nil
}

type Repository struct {
	Owner string
	Name  string
}

func (repo *Repository) Get(githubClient *github.Client) (*github.Repository, error) {
	repository, _, err := githubClient.Repositories.Get(repo.Owner, repo.Name)
	return repository, err
}

func (repo *Repository) CompareCommits(githubClient *github.Client, base, head string) (*github.CommitsComparison, error) {
	commitsComparison, _, err := githubClient.Repositories.CompareCommits(repo.Owner, repo.Name, base, head)
	return commitsComparison, err
}

type Target interface{}
