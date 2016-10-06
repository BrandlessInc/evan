package context

import (
	"fmt"

	"github.com/Everlane/evan/config"

	"github.com/google/go-github/github"
)

// Stores state relating to a deployment.
type Deployment struct {
	Application  *config.Application
	Environment  string
	Strategy     *config.Strategy
	Ref          string
	GithubClient *github.Client
	Initiator    interface{}
}

func (deployment *Deployment) GetApplication() *config.Application {
	return deployment.Application
}

func (deployment *Deployment) GetGithubClient() *github.Client {
	return deployment.GithubClient
}

func (deployment *Deployment) GetRef() string {
	return deployment.Ref
}

func (deployment *Deployment) GetInitiator() interface{} {
	return deployment.Initiator
}

func (deployment *Deployment) RunPreconditions() error {
	preconditions := deployment.Strategy.Preconditions

	resultChan := make(chan config.PreconditionResult)
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

	phases := deployment.Strategy.Phases
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

		switch status {
		case config.DONE:
			continue
		case config.IN_PROGRESS:
			// This "run" of the strategy is done for now if we're executing
			return nil
		case config.ERROR:
			if err != nil {
				return err
			} else {
				return fmt.Errorf("An unknown error occurred in phase: %v", phase)
			}
		default:
			return fmt.Errorf("Unknown status: %#v", status)
		}
	}

	return nil
}

func (deployment *Deployment) RunPhasePreloads() error {
	preloadablePhases := make([]config.PreloadablePhase, 0)
	for _, phase := range deployment.Strategy.Phases {
		if phase.CanPreload() {
			preloadablePhases = append(preloadablePhases, phase.(config.PreloadablePhase))
		}
	}

	resultChan := make(chan config.PreloadResult)
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
