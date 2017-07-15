package context

import (
	"fmt"

	"github.com/BrandlessInc/evan/common"

	"github.com/satori/go.uuid"
)

// Stores state relating to a deployment.
type Deployment struct {
	uuid        uuid.UUID
	application common.Application
	environment string
	strategy    common.Strategy
	ref         string
	sha1        string
	flags       map[string]interface{}
	products    map[string]interface{}

	store common.Store

	// Internal state
	currentState common.DeploymentState
	currentPhase common.Phase
	lastError    error
}

// Create a deployment for the given application to an environment.
func NewDeployment(app common.Application, environment string, strategy common.Strategy, ref string, flags map[string]interface{}) *Deployment {
	return &Deployment{
		uuid:         uuid.NewV1(),
		application:  app,
		environment:  environment,
		strategy:     strategy,
		ref:          ref,
		flags:        flags,
		products:     make(map[string]interface{}),
		currentState: common.DEPLOYMENT_PENDING,
	}
}

func NewBareDeployment() *Deployment {
	return &Deployment{
		flags:    make(map[string]interface{}),
		products: make(map[string]interface{}),
	}
}

func (deployment *Deployment) UUID() uuid.UUID {
	return deployment.uuid
}

func (deployment *Deployment) Application() common.Application {
	return deployment.application
}

func (deployment *Deployment) Environment() string {
	return deployment.environment
}

func (deployment *Deployment) Strategy() common.Strategy {
	return deployment.strategy
}

func (deployment *Deployment) Ref() string {
	return deployment.ref
}

func (deployment *Deployment) SHA1() string {
	return deployment.sha1
}

func (deployment *Deployment) SetSHA1(sha1 string) {
	deployment.sha1 = sha1
}

func (deployment *Deployment) MostPreciseRef() string {
	if deployment.sha1 != "" {
		return deployment.sha1
	} else {
		return deployment.ref
	}
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

func (deployment *Deployment) HasFlag(key string) bool {
	_, present := deployment.flags[key]
	return present
}

func (deployment *Deployment) Flag(key string) interface{} {
	return deployment.flags[key]
}

func (deployment *Deployment) SetFlag(key string, value interface{}) {
	deployment.flags[key] = value
}

func (deployment *Deployment) Products() map[string]interface{} {
	return deployment.products
}

func (deployment *Deployment) HasProduct(key string) bool {
	_, present := deployment.products[key]
	return present
}

func (deployment *Deployment) Product(key string) interface{} {
	return deployment.products[key]
}

func (deployment *Deployment) SetProduct(key string, value interface{}) {
	deployment.products[key] = value
}

// Looks for the "force" boolean in the `flags`.
func (deployment *Deployment) IsForce() bool {
	if force, ok := deployment.Flag("force").(bool); ok {
		return force
	} else {
		return false
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
	go deployment.notifyBeforePreconditions()
	defer func() { go deployment.notifyAfterPreconditions() }()

	deployment.setStateAndSave(common.RUNNING_PRECONDITIONS)

	preconditions := deployment.strategy.Preconditions()
	for _, precondition := range preconditions {
		err := precondition.Status(deployment)
		if err != nil {
			return err
		}
	}

	return nil
}

// Internal implementation of running phases. Manages setting
// `deployment.currentPhase` to the phase currently executing.
func (deployment *Deployment) runPhases(preloadResults PreloadResults) error {
	phases := deployment.strategy.Phases()
	for _, phase := range phases {
		deployment.currentPhase = phase

		preloadResult := preloadResults.Get(phase)
		err := phase.Execute(deployment, preloadResult)
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
	go deployment.notifyBeforePhases()
	defer func() { go deployment.notifyAfterPhases() }()

	results, err := deployment.RunPhasePreloads()
	if err != nil {
		deployment.lastError = err
		deployment.setStateAndSave(common.DEPLOYMENT_ERROR)
		return err
	}

	deployment.setStateAndSave(common.RUNNING_PHASE)

	err = deployment.runPhases(results)
	if err != nil {
		deployment.lastError = err
		deployment.setStateAndSave(common.DEPLOYMENT_ERROR)
		return err
	} else {
		deployment.setStateAndSave(common.DEPLOYMENT_DONE)
		return nil
	}
}

func (deployment *Deployment) callNotifiers(eachNotifier func(notifier common.Notifier, deployment *Deployment) error) {
	for _, notifier := range deployment.Strategy().Notifiers() {
		err := eachNotifier(notifier, deployment)
		if err != nil {
			fmt.Printf("Error notifying %v: %v\n", notifier, err)
		}
	}
}

func (deployment *Deployment) notifyBeforePreconditions() {
	deployment.callNotifiers(func(notifier common.Notifier, deployment *Deployment) error {
		return notifier.BeforePreconditions(deployment)
	})
}

func (deployment *Deployment) notifyAfterPreconditions() {
	deployment.callNotifiers(func(notifier common.Notifier, deployment *Deployment) error {
		return notifier.AfterPreconditions(deployment)
	})
}

func (deployment *Deployment) notifyBeforePhases() {
	deployment.callNotifiers(func(notifier common.Notifier, deployment *Deployment) error {
		return notifier.BeforePhases(deployment)
	})
}

func (deployment *Deployment) notifyAfterPhases() {
	deployment.callNotifiers(func(notifier common.Notifier, deployment *Deployment) error {
		return notifier.AfterPhases(deployment)
	})
}

type preloadResult struct {
	data interface{}
	err  error
}

type PreloadResults map[common.Phase]interface{}

func (results PreloadResults) Get(phase common.Phase) interface{} {
	return results[phase]
}

func (results PreloadResults) Set(phase common.Phase, data interface{}) {
	results[phase] = data
}

// Phases can expose preloads to gather any additional information they may
// need before executing. This will run those preloads in parallel.
func (deployment *Deployment) RunPhasePreloads() (PreloadResults, error) {
	preloadablePhases := make([]common.PreloadablePhase, 0)
	for _, phase := range deployment.strategy.Phases() {
		if phase.CanPreload() {
			preloadablePhases = append(preloadablePhases, phase.(common.PreloadablePhase))
		}
	}

	resultChan := make(chan preloadResult)
	for _, phase := range preloadablePhases {
		go func(currentPhase common.PreloadablePhase) {
			data, err := currentPhase.Preload(deployment)
			resultChan <- preloadResult{data: data, err: err}
		}(phase)
	}

	results := make(PreloadResults)
	for _, phase := range preloadablePhases {
		result := <-resultChan
		if result.err != nil {
			return nil, result.err
		} else {
			results.Set(phase.(common.Phase), result.data)
		}
	}

	return results, nil
}
