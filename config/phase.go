package config

type ExecuteStatus int

const (
	DONE ExecuteStatus = iota
	IN_PROGRESS
	ERROR
)

type Phase interface {
	CanPreload() bool

    // Has the phase already executed (synchronous, lifecycle).
    HasExecuted(Deployment) (bool, error)

    Execute(Deployment) (ExecuteStatus, error)
}

type PreloadResult struct {
	Data  interface{} // Data specific to the phase for it to use later on.
	Error error
}

type PreloadablePhase interface {
	Preload(Deployment) PreloadResult
}
