package preconditions

import (
	"time"

	"github.com/BrandlessInc/evan/common"
)

type SleepPrecondition struct {
	Duration time.Duration
}

func (s *SleepPrecondition) Status(deployment common.Deployment) error {
	time.Sleep(s.Duration)
	return nil
}
