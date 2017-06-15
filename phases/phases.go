// Contains all the built-in phases to be used for building deployment
// strategies.
package phases

import (
	"github.com/BrandlessInc/evan/common"
)

// Compiler verification that the implementations conform to the interface.
func _verify() []common.Phase {
	phases := make([]common.Phase, 0)
	phases = append(phases, &HerokuBuildPhase{})
	phases = append(phases, &SlackNotifierPhase{})
	return phases
}
