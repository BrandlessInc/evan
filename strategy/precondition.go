package strategy

import (
	"fmt"
)

type PreconditionResult struct {
	Precondition Precondition
	Error        error
}

type PreconditionResults chan PreconditionResult

type Precondition interface {
	Status(*Runner, PreconditionResults)
}

func createResult(precondition Precondition, err error) PreconditionResult {
	return PreconditionResult{
		Precondition: precondition,
		Error:        err,
	}
}

type GithubStatusesPrecondition struct{}

func (gh *GithubStatusesPrecondition) Status(runner *Runner, results PreconditionResults) {
	status, err := runner.Repository.GetGithubStatus(runner.GithubClient, runner.Ref)
	if err != nil {
		results <- createResult(gh, err)
		return
	}
	runner.CombinedStatus = status

	var result error = nil
	if *status.State != "success" {
		result = fmt.Errorf("Non-success status for ref: %v", *status.State)
	}
	results <- createResult(gh, result)
}
