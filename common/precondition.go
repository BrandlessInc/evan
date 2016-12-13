package common

type Precondition interface {
	Status(Deployment) error
}
