package common

type Strategy interface {
	Preconditions() []Precondition
	Phases() []Phase
	OnError(Deployment, error)
}
