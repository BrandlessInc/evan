// Built-in preconditions for use in deployment strategies.
package preconditions

import (
	"github.com/Everlane/evan/common"
)

// Compiler verification that the implementations conform to the interface.
func _verify() []common.Precondition {
	preconditions := make([]common.Precondition, 0)
	preconditions = append(preconditions, &GithubCombinedStatusPrecondition{})
	preconditions = append(preconditions, &GithubRequireAheadPrecondition{})
	return preconditions
}
