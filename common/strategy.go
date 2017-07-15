package common

type Strategy interface {
	Preconditions() []Precondition
	Phases() []Phase
	Notifiers() []Notifier
	OnError(Deployment, error)
}
