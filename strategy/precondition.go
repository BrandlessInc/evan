package strategy

type PreconditionResult struct {
    Precondition Precondition
    Error error
}

type PreconditionResults chan PreconditionResult

type Precondition interface {
	Status(*Strategy, PreconditionResults)
}

func createResult(precondition Precondition, err error) PreconditionResult {
    return PreconditionResult {
        Precondition: precondition,
        Error: err,
    }
}

type GithubStatusesPrecondition struct {}

func (gh *GithubStatusesPrecondition) Status(strategy *Strategy, results PreconditionResults) {
    status, err := strategy.Application.GetGithubStatus(strategy.Ref())
    if err != nil {
        results <- createResult(gh, err)
        return
    }
    strategy.SetGithubDeploymentStatus(status)

    // TODO: Actually verify that the combined status is okay.

    results <- createResult(gh, nil)
}
