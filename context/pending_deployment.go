package context

import (
	"github.com/Everlane/evan/common"
)

// Represents a deployment that wasn't able to be immediately deployed but
// should be deployed after a bit of time once its pending preconditions
// have been met.
type PendingDeployment struct {
	deployment *Deployment

	// List of preconditions which are still reporting a pending status.
	pending []common.Precondition
	// Preconditions that ended up failing.
	passed []common.Precondition
	// Preconditions that ended up passing.
	failed []common.Precondition
}

func NewPendingDeployment(deployment *Deployment, pending []common.Precondition) *PendingDeployment {
	return &PendingDeployment{
		deployment: deployment,
		pending:    pending,
		passed:     []common.Precondition{},
		failed:     []common.Precondition{},
	}
}
