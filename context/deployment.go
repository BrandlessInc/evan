package context

import (
	"fmt"

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
	initiator    interface{}
}

func NewDeployment(app *config.Application, environment string, strategy *config.Strategy, ref string) *Deployment {
	return &Deployment{
		application: app,
		environment: environment,
		strategy: strategy,
		ref: ref,
	}
}

func (deployment *Deployment) Application() common.Application {
	return deployment.application
}

func (deployment *Deployment) Ref() string {
	return deployment.ref
}

func (deployment *Deployment) GithubClient() *github.Client {
	return deployment.githubClient
}

func (deployment *Deployment) Initiator() interface{} {
	return deployment.initiator
}

func (deployment *Deployment) SetInitiator(initiator interface{}) {
	deployment.initiator = initiator
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

func (deployment *Deployment) RunPreconditions() error {
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

func (deployment *Deployment) RunPhases() error {
	err := deployment.RunPhasePreloads()
	if err != nil {
		return err
	}

	phases := deployment.strategy.Phases
	for _, phase := range phases {
		// Skip already-executed phases
		hasExecuted, err := phase.HasExecuted(deployment)
		if err != nil {
			return err
		}
		if hasExecuted {
			continue
		}

		status, err := phase.Execute(deployment)
		if err != nil {
			return err
		}

		switch status {
		case common.PHASE_DONE:
			continue
		case common.PHASE_IN_PROGRESS:
			// This "run" of the strategy is done for now if we're executing
			return nil
		case common.PHASE_ERROR:
			// We've already returned the error if it's present; so if we
			// reach here then it's `nil` and we don't know what's gone wrong
			return fmt.Errorf("An unknown error occurred in phase: %v", phase)
		default:
			return fmt.Errorf("Unknown status: %#v", status)
		}
	}

	return nil
}

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

func (deployment *Deployment) Run() error {
	err := deployment.RunPreconditions()
	if err != nil {
		return err
	}

	err = deployment.RunPhases()
	if err != nil {
		return err
	}

	return nil
}
