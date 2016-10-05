package config

type ExecuteStatus int

const (
	DONE ExecuteStatus = iota
	IN_PROGRESS
	ERROR
)

type Phase interface {
	CanPreload() bool

    // Returns whether or not the phase has executed (synchronous, lifecycle).
    CanExecute() (bool, error)
    // Is the phase currently executing (synchronous, lifecycle).
    IsExecuting() (bool, error)
    // Has the phase already executed (synchronous, lifecycle).
    HasExecuted() (bool, error)

    Execute(Deployment) (ExecuteStatus, error)
}

type PreloadResult struct {
	Phase Phase
	Data  interface{} // Data specific to the phase for it to use later on.
	Error error
}

type PreloadingPhase interface {
	Preload()
}
