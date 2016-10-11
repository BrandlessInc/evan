package config

import (
	"github.com/Everlane/evan/common"

	"github.com/google/go-github/github"
)

type Applications struct {
	List []*Application
	Store common.Store
}

func (apps *Applications) FindApplicationForGithubRepository(githubRepo *github.Repository) *Application {
	for _, app := range apps.List {
		appRepo := app.Repository

		if appRepo.Owner() == *githubRepo.Owner.Login && appRepo.Name() == *githubRepo.Name {
			return app
		}
	}

	return nil
}

func (apps *Applications) FindApplicationByName(name string) *Application {
	for _, app := range apps.List {
		// Ignore apps without a name
		if app.Name == "" {
			continue
		}

		if app.Name == name {
			return app
		}
	}

	return nil
}

type Application struct {
	Name       string
	Repository common.Repository

	Environments []string

	// Called to determine the target and strategy to use for deploying to a
	// given environment.
	DeployEnvironment func(string) *Strategy
}

func (app *Application) Wrapper() *CommonApplicationWrapper {
	return &CommonApplicationWrapper{
		app: app,
	}
}

// Wraps `Application` struct to fulfill the `common.Application` interface.
type CommonApplicationWrapper struct {
	app *Application
}

func (wrapper *CommonApplicationWrapper) Name() string {
	return wrapper.app.Name
}

func (wrapper *CommonApplicationWrapper) Repository() common.Repository {
	return wrapper.app.Repository
}

func (wrapper *CommonApplicationWrapper) StrategyForEnvironment(environment string) common.Strategy {
	environmentFound := false
	for _, configuredEnvironment := range wrapper.app.Environments {
		if environment == configuredEnvironment {
			environmentFound = true
			break
		}
	}
	if !environmentFound {
		return nil
	}

	strategy := wrapper.app.DeployEnvironment(environment)
	if strategy == nil {
		return nil
	}

	return &CommonStrategyWrapper{strategy: strategy}
}

type Repository struct {
	owner string
	name  string
}

func NewRepository(owner, name string) *Repository {
	return &Repository{owner: owner, name: name}
}

func (repo *Repository) Owner() string {
	return repo.owner
}

func (repo *Repository) Name() string {
	return repo.name
}

type Target interface{}
