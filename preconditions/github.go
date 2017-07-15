package preconditions

import (
	"fmt"

	"github.com/BrandlessInc/evan/common"

	"github.com/google/go-github/github"
)

// Fetches the SHA1 for the commit from GitHub.
type GithubFetchCommitSHA1Precondition struct{}

func (gh *GithubFetchCommitSHA1Precondition) Status(deployment common.Deployment) error {
	githubRepo, err := common.NewGithubRepositoryFromDeployment(deployment)
	if err != nil {
		return err
	}

	sha1, err := githubRepo.GetCommitSHA1(deployment.Ref())
	if err != nil {
		return err
	}
	deployment.SetSHA1(sha1)

	return nil
}

type GithubCombinedStatusPrecondition struct {
	// If true then it will ignore the reported status if GitHub reports that
	// there are no status checks.
	AllowEmpty bool
}

func (gh *GithubCombinedStatusPrecondition) Status(deployment common.Deployment) error {
	repo := deployment.Application().Repository()
	ref := deployment.MostPreciseRef()

	client, err := deployment.GithubClient()
	if err != nil {
		return err
	}

	status, _, err := client.Repositories.GetCombinedStatus(repo.Owner(), repo.Name(), ref, nil)
	if err != nil {
		return err
	}

	// Skip if there are no status checks and that's allowed
	if *status.TotalCount == 0 && gh.AllowEmpty {
		return nil
	}

	switch *status.State {
	case "success":
		return nil
	case "pending":
		return common.NewPendingError(fmt.Errorf("Status pending"))
	default:
		return fmt.Errorf("Non-success status for ref: %v", *status.State)
	}
}

// Require that the branch for deployment not be behind the default branch
// on GitHub. If `AutoMerge` is true then it will try to create a merge via
// the GitHub API if the deployed ref is behind.
type GithubRequireAheadPrecondition struct {
	AutoMerge bool
}

type GithubRequireAheadContext struct {
	RepoClient  *common.GithubRepository
	RepoDetails *github.Repository
}

// Compares the ref being deployed against the default branch on GitHub to
// determine whether or not a merge needs to happen. Returns `false` if it's
// a force deployment.
func (gh *GithubRequireAheadPrecondition) NeedsMerge(deployment common.Deployment, ctx *GithubRequireAheadContext) (bool, error) {
	if deployment.IsForce() {
		return false, nil
	}

	base := *ctx.RepoDetails.DefaultBranch
	head := deployment.Ref()

	comparison, err := ctx.RepoClient.CompareCommits(base, head)
	if err != nil {
		return false, err
	}

	return (*comparison.BehindBy > 0), nil
}

// Creates merge commit to get the target branch (deployment's ref) up-to-date
// with the default branch of the repository.
func (gh *GithubRequireAheadPrecondition) Merge(deployment common.Deployment, ctx *GithubRequireAheadContext) (string, error) {
	base := deployment.Ref()
	head := *ctx.RepoDetails.DefaultBranch
	commitMessage := fmt.Sprintf("Merge '%v' into '%v'", head, base)
	commit, err := ctx.RepoClient.Merge(base, head, commitMessage)
	if err != nil {
		return "", err
	} else {
		return *commit.SHA, nil
	}
}

func (gh *GithubRequireAheadPrecondition) Status(deployment common.Deployment) error {
	githubClient, err := deployment.GithubClient()
	if err != nil {
		return err
	}

	repoClient := &common.GithubRepository{
		Repository:   deployment.Application().Repository(),
		GithubClient: githubClient,
	}
	repoDetails, err := repoClient.Get()
	if err != nil {
		return err
	}

	ctx := &GithubRequireAheadContext{
		RepoClient:  repoClient,
		RepoDetails: repoDetails,
	}

	needsMerge, err := gh.NeedsMerge(deployment, ctx)
	if err != nil {
		return err
	}

	// Halt if we don't need to merge!
	if !needsMerge {
		return nil
	}

	if !gh.AutoMerge {
		return fmt.Errorf("Merge needed for ref '%v'", deployment.Ref())
	}

	sha1, err := gh.Merge(deployment, ctx)
	if err != nil {
		return err
	}
	// Update the SHA1 to point to the new merge commit
	deployment.SetSHA1(sha1)

	return nil
}
