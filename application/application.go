package application

import (
	"github.com/Everlane/evan/repository"
	"github.com/Everlane/evan/strategy"
)

// A single code-base deployed to 1+ targets for 1+ environments.
type Application struct {
	Targets                map[string]Target
	Environments           []string
	TargetForEnvironment   func(string) *Target
	StrategyForEnvironment func(string) *strategy.Strategy

	// Details of the GitHub repository corresponding to this code-base.
	Repository *repository.Repository
}

type Target interface{}
