package config

type PreconditionResult struct {
	Precondition Precondition
	Error        error
}

func createResult(precondition Precondition, err error) PreconditionResult {
	return PreconditionResult{
		Precondition: precondition,
		Error:        err,
	}
}

type Precondition interface {
	Status(Deployment) PreconditionResult
}
