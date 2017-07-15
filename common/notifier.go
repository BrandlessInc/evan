package common

type Notifier interface {
	BeforePreconditions(Deployment) error
	AfterPreconditions(Deployment) error
	BeforePhases(Deployment) error
	AfterPhases(Deployment) error
}
