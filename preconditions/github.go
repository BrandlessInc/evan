package preconditions

import (
	"fmt"

	"github.com/Everlane/evan/common"
)

type GithubCombinedStatusPrecondition struct{}

func (gh *GithubCombinedStatusPrecondition) Status(deployment common.Deployment) common.PreconditionResult {
	repo := deployment.Application().Repository()
	ref := deployment.Ref()

	client, err := common.GithubClient(deployment)
	if err != nil {
		return createResult(gh, err)
	}

	status, _, err := client.Repositories.GetCombinedStatus(repo.Owner(), repo.Name(), ref, nil)
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
func (gh *GithubRequireAheadPrecondition) NeedsMerge(deployment common.Deployment) (bool, error) {
	if deployment.IsForce() {
		return false, nil
	}

	githubClient, err := common.GithubClient(deployment)
	if err != nil {
		return false, err
	}

	githubRepo := &common.GithubRepository{
		Repository:   deployment.Application().Repository(),
		GithubClient: githubClient,
	}

	repoDetails, err := githubRepo.Get()
	if err != nil {
		return false, err
	}

	base := *repoDetails.DefaultBranch
	head := deployment.Ref()

	comparison, err := githubRepo.CompareCommits(base, head)
	if err != nil {
		return false, err
	}

	return (*comparison.BehindBy > 0), nil
}

func (gh *GithubRequireAheadPrecondition) Status(deployment common.Deployment) common.PreconditionResult {
	return createResult(gh, nil)
}
