package config

import (
	"github.com/google/go-github/github"
)

type Deployment interface {
	GetApplication() *Application
	GetRef() string
	GetGithubClient() *github.Client
}
