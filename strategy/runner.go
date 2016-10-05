package strategy

import (
	"github.com/Everlane/evan/repository"

	"github.com/google/go-github/github"
)

// Represents the state of a strategy as it is being run.
type Runner struct {
	// External configuration
	Repository   *repository.Repository
	GithubClient *github.Client

	// Internal configuration
	Strategy *Strategy

	// Git ref for which we're running the strategy.
	Ref string
	// Result of the combined status on GitHub for the ref.
	CombinedStatus *github.CombinedStatus
}
