package config

import (
	"fmt"
)

type GithubCombinedStatusPrecondition struct{}

func (gh *GithubCombinedStatusPrecondition) Status(app *Application, deployment Deployment, results PreconditionResults) {
	repo := app.Repository
	ref := deployment.GetRef()
	client := deployment.GetGithubClient()

	status, _, err := client.Repositories.GetCombinedStatus(repo.Owner, repo.Name, ref, nil)
	if err != nil {
		results <- createResult(gh, err)
		return
	}

	var result error = nil
	if *status.State != "success" {
		result = fmt.Errorf("Non-success status for ref: %v", *status.State)
	}
	results <- createResult(gh, result)
}
