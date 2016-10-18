package common

import (
	"github.com/satori/go.uuid"
)

type DeploymentState int

const (
	DEPLOYMENT_PENDING DeploymentState = iota
	RUNNING_PRECONDITIONS
	RUNNING_PHASE
	DEPLOYMENT_DONE
	DEPLOYMENT_ERROR
)

func (state DeploymentState) String() string {
	switch state {
	case DEPLOYMENT_PENDING:
		return "DEPLOYMENT_PENDING"
	case RUNNING_PRECONDITIONS:
		return "RUNNING_PRECONDITIONS"
	case RUNNING_PHASE:
		return "RUNNING_PHASE"
	case DEPLOYMENT_DONE:
		return "DEPLOYMENT_DONE"
	case DEPLOYMENT_ERROR:
		return "DEPLOYMENT_ERROR"
	default:
		return "UNKNOWN"
	}
}

type DeploymentStatus struct {
	State DeploymentState
	Phase Phase
	Error error
}

type Deployment interface {
	UUID() uuid.UUID
	Application() Application
	Environment() string
	Strategy() Strategy
	Ref() string
	SHA1() string
	SetSHA1(string)
	// Return the SHA1 if known, otherwise the ref
	MostPreciseRef() string
	Flags() map[string]interface{}
	HasFlag(string) bool
	Flag(string) interface{}
	SetFlag(string, interface{})
	IsForce() bool

	Status() DeploymentStatus
}

type Strategy interface {
	Preconditions() []Precondition
	Phases() []Phase
	OnError(Deployment, error)
}
