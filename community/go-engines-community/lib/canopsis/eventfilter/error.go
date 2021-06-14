package eventfilter

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
)

// A DropError is an error that is returned by the eventfilter when an event
// should be dropped.
type DropError interface {
	errt.ErrT
	IsDropError()
}

type dropError struct {
	errt.ErrT
}

// IsDropError is a function that does nothing. It is here to differentiate the
// DropError and ErrT interfaces.
func (e dropError) IsDropError() {}

// NewDropError wraps an error into a DropError.
func NewDropError(err error) DropError {
	if err == nil {
		return nil
	}
	return dropError{
		ErrT: errt.NewErrT(err),
	}
}

// DefaultDropError creates a new DropError with a default message.
func DefaultDropError() DropError {
	return NewDropError(fmt.Errorf("dropping event"))
}
