package common

type Phase interface {
	CanPreload() bool
	// Second argument is the result from the preload, or `nil` if the phase
	// doesn't preload.
	Execute(Deployment, interface{}) error
}

type PreloadablePhase interface {
	// Returns data specific to the phase for it to use later on; if it
	// returns an error then the deployment is cancelled.
	Preload(Deployment) (interface{}, error)
}
