package config

import (
	"github.com/Everlane/evan/common"

	"github.com/google/go-github/github"
)

type Applications struct {
	Map map[string]*Application
}

func (apps *Applications) FindApplicationForGithubRepository(githubRepo *github.Repository) *Application {
	for _, app := range apps.Map {
		appRepo := app.Repository

		if appRepo.Owner() == *githubRepo.Owner.Login && appRepo.Name() == *githubRepo.Name {
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
	strategy := wrapper.app.DeployEnvironment(environment)
	if strategy == nil {
		return nil
	}

	return &CommonStrategyWrapper{strategy: strategy}
}

type Target interface{}
