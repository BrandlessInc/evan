package common

type Store interface {
    SaveDeployment(Deployment) error
    // Look up a `Deployment` by its environment and ref.
    FindDeployment(environment string, ref string) (Deployment, error)

    ShouldCancel() (bool, error)
    SetShouldCancel(bool) error
}
