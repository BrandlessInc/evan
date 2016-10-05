package evan

type StateKey string

const (
	GITHUB_DEPLOYMENT_EVENT = "github.deployment_event"
)

type StateMap map[StateKey]interface{}

// Describes how an application will be deployed to a environment & target.
type Strategy struct {
	// State is initialized when the strategy is executed and built up as the
	// preconditions and reporters run.
	State StateMap

	Preconditions []Precondition
	Phases        []Phase
	Reporter      Reporter
}

func newStrategyWithDefaults() *Strategy {
	return &Strategy{
		State:         make(StateMap),
		Preconditions: make([]Precondition, 0),
		Phases:        make([]Phase, 0),
		Reporter:      nil,
	}
}

func NewStrategyFromGithubDeploymentEvent(deploymentEvent interface{}) *Strategy {
	strategy := newStrategyWithDefaults()
	strategy.State[GITHUB_DEPLOYMENT_EVENT] = deploymentEvent
	return strategy
}

type Precondition interface {}
type Phase interface {}
type Reporter interface {}
