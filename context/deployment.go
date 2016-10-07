package context

import (
	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/config"

	"github.com/google/go-github/github"
)

// Stores state relating to a deployment.
type Deployment struct {
	application  *config.Application
	environment  string
	strategy     *config.Strategy
	ref          string
	flags        map[string]interface{}

	githubClient *github.Client
	store common.Store

	// Internal state
	currentState common.DeploymentState
	currentPhase common.Phase
	lastError error
}

func NewDeployment(app *config.Application, environment string, strategy *config.Strategy, ref string) *Deployment {
	return &Deployment{
		application: app,
		environment: environment,
		strategy: strategy,
		ref: ref,
		currentState: common.DEPLOYMENT_PENDING,
	}
}

func (deployment *Deployment) Application() common.Application {
	return deployment.application
}

func (deployment *Deployment) Environment() string {
	return deployment.environment
}

func (deployment *Deployment) Ref() string {
	return deployment.ref
}

func (deployment *Deployment) GithubClient() *github.Client {
	return deployment.githubClient
}

func (deployment *Deployment) SetGithubClient(githubClient *github.Client) {
	deployment.githubClient = githubClient
}

func (deployment *Deployment) SetStoreAndSave(store common.Store) error {
	deployment.store = store
	return store.SaveDeployment(deployment)
}

// Will panic if it is unable to save. This will be called *after*
// `SetStoreAndSave` should have been called, so we're assuming that if that
// worked then this should also work.
func (deployment *Deployment) setStateAndSave(state common.DeploymentState) {
	deployment.currentState = state
	err := deployment.store.SaveDeployment(deployment)
	if err != nil {
		panic(err)
	}
}

func (deployment *Deployment) Flags() map[string]interface{} {
	return deployment.flags
}

// Looks for the "force" boolean in the `flags`.
func (deployment *Deployment) IsForce() bool {
	forceUntyped, present := deployment.flags["force"]
	if !present {
		return false
	}

	force, ok := forceUntyped.(bool)
	if !ok || !force {
		return false
	} else {
		return true
	}
}

func (deployment *Deployment) Status() common.DeploymentStatus {
	var phase common.Phase
	if deployment.currentState == common.RUNNING_PHASE {
		phase = deployment.currentPhase
	}

	return common.DeploymentStatus{
		State: deployment.currentState,
		Phase: phase,
		Error: nil,
	}
}

func (deployment *Deployment) CheckPreconditions() error {
	deployment.setStateAndSave(common.RUNNING_PRECONDITIONS)

	preconditions := deployment.strategy.Preconditions

	resultChan := make(chan common.PreconditionResult)
	for _, precondition := range preconditions {
		go func() {
			resultChan <- precondition.Status(deployment)
		}()
	}

	for _ = range preconditions {
		result := <-resultChan
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Internal implementation of running phases. Manages setting
// `deployment.currentPhase` to the phase currently executing.
func (deployment *Deployment) runPhases() error {
	phases := deployment.strategy.Phases
	for _, phase := range phases {
		deployment.currentPhase = phase

		err := phase.Execute(deployment)
		if err != nil {
			return err
		}
	}

	return nil
}

// Runs all the phases configured in the `Strategy`. Sets `currentState` and
// `currentPhase` fields as appropriate. If an error occurs it will also set
// the `lastError` field to that error.
func (deployment *Deployment) RunPhases() error {
	err := deployment.RunPhasePreloads()
	if err != nil {
		return err
	}

	deployment.setStateAndSave(common.RUNNING_PHASE)

	err = deployment.runPhases()
	if err != nil {
		deployment.lastError = err
		deployment.setStateAndSave(common.DEPLOYMENT_ERROR)
		return err
	} else {
		deployment.setStateAndSave(common.DEPLOYMENT_DONE)
		return nil
	}
}

// Phases can expose preloads to gather any additional information they may
// need before executing. This will run those preloads in parallel.
func (deployment *Deployment) RunPhasePreloads() error {
	preloadablePhases := make([]common.PreloadablePhase, 0)
	for _, phase := range deployment.strategy.Phases {
		if phase.CanPreload() {
			preloadablePhases = append(preloadablePhases, phase.(common.PreloadablePhase))
		}
	}

	resultChan := make(chan common.PreloadResult)
	for _, phase := range preloadablePhases {
		go func() {
			resultChan <- phase.Preload(deployment)
		}()
	}

	for _ = range preloadablePhases {
		result := <-resultChan
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
