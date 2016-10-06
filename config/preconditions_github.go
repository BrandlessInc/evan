package config

import (
	"fmt"
)

type GithubCombinedStatusPrecondition struct{}

func (gh *GithubCombinedStatusPrecondition) Status(deployment Deployment) PreconditionResult {
	repo := deployment.GetApplication().Repository
	ref := deployment.GetRef()
	client := deployment.GetGithubClient()

	status, _, err := client.Repositories.GetCombinedStatus(repo.Owner, repo.Name, ref, nil)
	if err != nil {
		return createResult(gh, err)
	}

	var result error = nil
	if *status.State != "success" {
		result = fmt.Errorf("Non-success status for ref: %v", *status.State)
	}
	return createResult(gh, result)
}
