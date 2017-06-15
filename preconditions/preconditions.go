// Built-in preconditions for use in deployment strategies.
package preconditions

import (
	"github.com/BrandlessInc/evan/common"
)

// Compiler verification that the implementations conform to the interface.
func _verify() []common.Precondition {
	return []common.Precondition{
		&GithubCombinedStatusPrecondition{},
		&GithubFetchCommitSHA1Precondition{},
		&GithubRequireAheadPrecondition{},
		&RestrictForcePrecondition{},
	}
}
