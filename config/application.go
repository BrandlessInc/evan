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

type Target interface{}
