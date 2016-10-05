package application

import (
	"github.com/google/go-github/github"
)

type Repository struct {
	Owner string
	Name  string
}

// A single code-base deployed to 1+ targets for 1+ environments.
type Application struct {
	Targets              map[string]Target
	Environments         []string
	TargetForEnvironment func(string) *Target

	// GitHub API client to use.
	GithubClient *github.Client
	// Details of the GitHub repository corresponding to this code-base.
	Repository Repository
}

func (app *Application) GetGithubStatus(ref string) (*github.CombinedStatus, error) {
	status, _, err := app.GithubClient.Repositories.GetCombinedStatus(app.Repository.Owner, app.Repository.Name, ref, nil)
	return status, err
}

type Target interface{}
