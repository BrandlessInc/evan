package common

type Precondition interface {
	Status(Deployment) error
}

// Wraps an error to signify that it's not a true error but rather pending
// resolution.
type PendingError struct {
	error
}

func NewPendingError(err error) *PendingError {
	return &PendingError{err}
}
