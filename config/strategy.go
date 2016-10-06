package config

import (
	"github.com/Everlane/evan/common"
)

type Strategy struct {
	Preconditions []common.Precondition
	Phases        []common.Phase
}
