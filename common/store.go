package common

type Store interface {
	SaveDeployment(Deployment) error
	// Look up a `Deployment` for an application environment.
	FindDeployment(application Application, environment string) (Deployment, error)

	// Enqueue a deployment to be executed immediately after the current
	// deployment finishes.
	EnqueueDeployment(Deployment) error
	FindEnqueuedDeployments(application Application, environment string) ([]Deployment, error)

	// ShouldCancel() (bool, error)
	// SetShouldCancel(bool) error

	// Returns whether or not there is an in-progress deployment to a given
	// application environment.
	HasActiveDeployment(Application, string) (bool, error)
}
