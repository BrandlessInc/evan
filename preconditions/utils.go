package preconditions

import (
    "github.com/Everlane/evan/common"
)

func createResult(precondition common.Precondition, err error) common.PreconditionResult {
	return common.PreconditionResult{
		Precondition: precondition,
		Error:        err,
	}
}
