package common

type PreconditionResult struct {
	Precondition Precondition
	Error        error
}

type Precondition interface {
	Status(Deployment) PreconditionResult
}
