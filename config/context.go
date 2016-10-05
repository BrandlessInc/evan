package config

import (
	"github.com/google/go-github/github"
)

type Deployment interface {
	GetRef() string
	GetGithubClient() *github.Client
}
