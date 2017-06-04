package config

import (
	"github.com/BrandlessInc/evan/common"
)

type Strategy struct {
	Preconditions []common.Precondition
	Phases        []common.Phase
	OnError       func(common.Deployment, error)
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

func (wrapper *CommonStrategyWrapper) OnError(deployment common.Deployment, err error) {
	if wrapper.strategy.OnError != nil {
		wrapper.strategy.OnError(deployment, err)
	}
}
