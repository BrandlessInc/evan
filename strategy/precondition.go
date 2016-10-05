package strategy

type PreconditionResult struct {
    Precondition Precondition
    Error error
}

type PreconditionResults chan PreconditionResult

type Precondition interface {
	Status(*Runner, PreconditionResults)
}

func createResult(precondition Precondition, err error) PreconditionResult {
    return PreconditionResult {
        Precondition: precondition,
        Error: err,
    }
}

type GithubStatusesPrecondition struct {}

func (gh *GithubStatusesPrecondition) Status(runner *Runner, results PreconditionResults) {
    status, err := runner.Application.GetGithubStatus(runner.Ref)
    if err != nil {
        results <- createResult(gh, err)
        return
    }
    runner.CombinedStatus = status

    // TODO: Actually verify that the combined status is okay.

    results <- createResult(gh, nil)
}
