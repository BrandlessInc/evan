package common

type ExecuteStatus int

const (
	PHASE_DONE ExecuteStatus = iota
	PHASE_IN_PROGRESS
	PHASE_ERROR
)

type Phase interface {
	CanPreload() bool
	Execute(Deployment) error
}

type PreloadResult struct {
	Data  interface{} // Data specific to the phase for it to use later on.
	Error error
}

type PreloadablePhase interface {
	Preload(Deployment) PreloadResult
}
