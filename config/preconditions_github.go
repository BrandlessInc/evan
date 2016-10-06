package config

import (
	"fmt"
)

type GithubCombinedStatusPrecondition struct{}

func (gh *GithubCombinedStatusPrecondition) Status(deployment Deployment) PreconditionResult {
	repo := deployment.Application().Repository
	ref := deployment.Ref()
	client := deployment.GithubClient()

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

// Require that the branch for deployment not be behind the default branch
// on GitHub. If `AutoMerge` is true then it will try to create a merge via
// the GitHub API if the deployed ref is behind.
type GithubRequireAheadPrecondition struct {
	AutoMerge bool
}

// Compares the ref being deployed against the default branch on GitHub to
// determine whether or not a merge needs to happen. Returns `false` if it's
// a force deployment.
func (gh *GithubRequireAheadPrecondition) NeedsMerge(deployment Deployment) (bool, error) {
	if deployment.IsForce() {
		return false, nil
	}

	githubClient := deployment.GithubClient()
	repo := deployment.Application().Repository

	repoOnGithub, err := repo.Get(githubClient)
	if err != nil {
		return false, err
	}

	base := *repoOnGithub.DefaultBranch
	head := deployment.Ref()

	comparison, err := repo.CompareCommits(githubClient, base, head)
	if err != nil {
		return false, err
	}

	return (*comparison.BehindBy > 0), nil
}
