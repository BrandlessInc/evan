package preconditions

import (
	"fmt"

	"github.com/Everlane/evan/common"
)

type RestrictForcePrecondition struct {
	// Array of environments for which force deploys are allowed.
	Safelist []string
}

func (rfp *RestrictForcePrecondition) Status(deployment common.Deployment) error {
	if !deployment.IsForce() {
		return nil
	}

	for _, safeEnvironment := range rfp.Safelist {
		if deployment.Environment() == safeEnvironment {
			return nil
		}
	}

	return fmt.Errorf("Cannot force deploy to %s", deployment.Environment())
}
