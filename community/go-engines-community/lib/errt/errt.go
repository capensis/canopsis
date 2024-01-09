// Package errt package defines some basic errors to use, see example.
//
// The idea is to encapsulate any error in a hidden struct that implements a
// public interface.
//
// That way, type assertion will give you the kind of error you encountered.
//
// Also, you're free to redefine your custom errors in interfaces to avoid
// tight coupling between your program/library and this one.
package errt

// ErrT is the generic interface that all errt-compatible errors must implement.
type ErrT interface {
	// Returns the original error
	Err() error
	error
}

// RefErrT implements ErrT interface and avoids you some heavy useless code duplication.
type refErrT struct {
	err error
}

func (e refErrT) Error() string {
	return e.err.Error()
}

func (e refErrT) Err() error {
	return e.err
}

// NewErrT creates a RefErrT
func NewErrT(err error) ErrT {
	return refErrT{
		err: err,
	}
}

type ioerror struct {
	ErrT
}

func (e ioerror) IsIOError() {
}

// IOError of any kind: timeout, host unreachable, file permissions...
type IOError interface {
	IsIOError()
	ErrT
}

// NewIOError returns an IO Error
func NewIOError(err error) error {
	if err == nil {
		return nil
	}
	return ioerror{ErrT: NewErrT(err)}
}

type notFound struct {
	ErrT
}

func (e notFound) IsNotFound() {
}

// NotFound error
type NotFound interface {
	IsNotFound()
	ErrT
}

// NewNotFound returns a NotFound Error
func NewNotFound(err error) error {
	if err == nil {
		return nil
	}
	return notFound{ErrT: NewErrT(err)}
}

// Duplicated can be used when something is a duplicate of something else.
type Duplicated interface {
	IsDuplicated()
	ErrT
}

type duplicated struct {
	ErrT
}

// NewDuplicated returns a Duplicated Error
func NewDuplicated(err error) error {
	if err == nil {
		return nil
	}
	return duplicated{ErrT: NewErrT(err)}
}

// Fatal error that should lead to program exit
type Fatal interface {
	IsFatal()
	RCode() int
	ErrT
}

type fatal struct {
	ErrT
	rcode int
}

// RCode is useful to embed a return code
func (e fatal) RCode() int {
	return e.rcode
}

func (e fatal) IsFatal() {
}

// NewFatal returns a new Fatal Error
//
// param rcode is the optional return code that can be used for os.Exit() for example.
func NewFatal(err error, rcode int) error {
	if err == nil {
		return nil
	}
	return fatal{
		ErrT:  NewErrT(err),
		rcode: rcode,
	}
}

// UnmanagedEventError can be returned from any method that doesn't manage a given types.Event
type UnmanagedEventError interface {
	ErrT
	IsUnmanagedEventError()
}

type unmanagedEventError struct {
	ErrT
}

func (e unmanagedEventError) IsUnmanagedEventError() {
}

// NewUnmanagedEventError ...
func NewUnmanagedEventError(err error) error {
	if err == nil {
		return nil
	}
	return unmanagedEventError{ErrT: NewErrT(err)}
}

// UnknownError should be used to encapsulate any non fatal error with
// no consequences other than "it doesn't worked"
type UnknownError interface {
	ErrT
	IsUnknownError()
}

type unknownError struct {
	ErrT
}

func (e unknownError) IsUnknownError() {
}

// NewUnknownError ...
func NewUnknownError(err error) error {
	if err == nil {
		return nil
	}

	return unknownError{ErrT: NewErrT(err)}
}
