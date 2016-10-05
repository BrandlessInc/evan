package application

import (
	"github.com/Everlane/evan/repository"
	"github.com/Everlane/evan/strategy"
	"github.com/Everlane/evan/target"

	"github.com/google/go-github/github"
)

// A single code-base deployed to 1+ targets for 1+ environments.
type Application struct {
	Targets                map[string]*target.Target
	Environments           []string
	TargetForEnvironment   func(string) *target.Target
	StrategyForEnvironment func(string) *strategy.Strategy

	// Details of the GitHub repository corresponding to this code-base.
	Repository *repository.Repository
}

type Suite struct {
	Applications []Application
}

func (suite *Suite) FindApplicationForGithubRepository(repo *github.Repository) *Application {
	name := *repo.Name
	owner := *repo.Owner.Login

	for _, app := range suite.Applications {
		if app.Repository.Owner == owner && app.Repository.Name == name {
			return &app
		}
	}

	return nil
}
