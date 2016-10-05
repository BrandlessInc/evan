package strategy

import (
    "github.com/google/go-github/github"

    "github.com/Everlane/evan/application"
)

// Represents the state of a strategy as it is being run.
type Runner struct {
    Application *application.Application
    Strategy *Strategy

    Ref string
    CombinedStatus *github.CombinedStatus
}
