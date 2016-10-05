package strategy

import (
	"github.com/google/go-github/github"

	"github.com/Everlane/evan/application"
)

type StateKey string

const (
	REF                      = "ref"
	GITHUB_DEPLOYMENT_STATUS = "github.deployment_status"
	GITHUB_DEPLOYMENT        = "github.deployment"
)

type StateMap map[StateKey]interface{}

// Describes how an application will be deployed to a environment & target.
type Strategy struct {
	// State is initialized when the strategy is executed and built up as the
	// preconditions and reporters run.
	State StateMap

	Application   *application.Application
	Preconditions []Precondition
	Phases        []Phase
	Reporter      Reporter
}

func newStrategyWithDefaults(Application *application.Application) *Strategy {
	return &Strategy{
		State:       make(StateMap),
		Reporter:    nil,
		Application: Application,
	}
}

// Returns the ref for which this deployment strategy is running.
func (strategy *Strategy) Ref() string {
	return strategy.State[REF].(string)
}

func (strategy *Strategy) SetGithubDeploymentStatus(status *github.CombinedStatus) {
	strategy.State[GITHUB_DEPLOYMENT_STATUS] = status
}

type Phase interface{}
type Reporter interface{}
