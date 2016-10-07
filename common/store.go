package common

type Store interface {
	SaveDeployment(Deployment) error
	// Look up a `Deployment` for an application to an environment.
	FindDeployment(application Application, environment string) (Deployment, error)

	// ShouldCancel() (bool, error)
	// SetShouldCancel(bool) error
}
