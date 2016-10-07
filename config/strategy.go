package config

import (
	"github.com/Everlane/evan/common"
)

type Strategy struct {
	Preconditions []common.Precondition
	Phases        []common.Phase
}

type CommonStrategyWrapper struct {
	strategy *Strategy
}

func (wrapper *CommonStrategyWrapper) Preconditions() []common.Precondition {
	return wrapper.strategy.Preconditions
}

func (wrapper *CommonStrategyWrapper) Phases() []common.Phase {
	return wrapper.strategy.Phases
}
