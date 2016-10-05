package repository

import (
    "github.com/google/go-github/github"
)

type Repository struct {
	Owner string
	Name  string
}

func (repo *Repository) GetGithubStatus(client *github.Client, ref string) (*github.CombinedStatus, error) {
	status, _, err := client.Repositories.GetCombinedStatus(repo.Owner, repo.Name, ref, nil)
	return status, err
}
