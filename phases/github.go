package phases

import (
	"fmt"
	"time"

	"github.com/BrandlessInc/evan/common"
)

// Poll the GitHub API to check the combined status for the deployment's
// Git ref/SHA1. By default this will check once; increase the Timeout
// amount to allow it to check repeatedly. To avoid flooding the Github API
// it will wait 10 seconds between checks.
type GithubCombinedStatusPhase struct {
	// Treat no status checks as okay (ignores status reported by GitHub)
	AllowEmpty bool
	// How long to wait for the combined status to no longer be pending
	Timeout time.Duration
}

func (gh *GithubCombinedStatusPhase) CanPreload() bool {
	return false
}

func (gh *GithubCombinedStatusPhase) Execute(deployment common.Deployment, data interface{}) error {
	repo := deployment.Application().Repository()
	ref := deployment.MostPreciseRef()

	client, err := deployment.GithubClient()
	if err != nil {
		return err
	}

	start := time.Now()

	for true {
		status, _, err := client.Repositories.GetCombinedStatus(repo.Owner(), repo.Name(), ref, nil)
		if err != nil {
			return err
		}

		// Skip if having no status checks is allowed
		if *status.TotalCount == 0 && gh.AllowEmpty {
			return nil
		}

		state := *status.State

		switch state {
		case "pending":
			elapsed := time.Since(start)
			if elapsed > gh.Timeout {
				return fmt.Errorf("Timed out waiting for combined status")
			}
			time.Sleep(10 * time.Second)
			continue
		case "success":
			return nil
		case "failure":
			return fmt.Errorf("Combined status is failure")
		default:
			return fmt.Errorf("Unknown combined status: %v", state)
		}
	}

	panic("Unreachable")
}
