package mongo

import (
	"io"

	"github.com/globalsign/mgo"

	"git.canopsis.net/canopsis/go-engines/lib/errt"
)

type MongoError interface {
	IsMongoError()
	errt.ErrT
}

type MongoConnectionError interface {
	MongoError
	IsMongoConnectionError()
}

type mongoError struct {
	errt.ErrT
}

// IsMongoError is a function that does nothing. It is here to differentiate
// the MongoError and ErrT interfaces.
func (e mongoError) IsMongoError() {
}

func NewMongoError(err error) error {
	if err == nil {
		return nil
	}
	return mongoError{ErrT: errt.NewErrT(err)}
}

type mongoConnectionError struct {
	errt.ErrT
}

func (e mongoConnectionError) IsMongoError() {
}

func (e mongoConnectionError) IsMongoConnectionError() {
}

// WrapError is a helper to convert an error to an errt error.
// If the error is unknown, returns MongoError.
func WrapError(err error) error {
	if err == nil {
		return nil
	}

	if mgo.IsDup(err) {
		return errt.NewDuplicated(err)
	}

	if err == mgo.ErrNotFound {
		return errt.NewNotFound(err)
	}

	if err.Error() == "no reachable servers" {
		return mongoConnectionError{ErrT: errt.NewErrT(err)}
	}

	if err == io.EOF {
		return errt.NewIOError(err)
	}

	if err == io.ErrClosedPipe {
		return errt.NewIOError(err)
	}

	if err == io.ErrNoProgress {
		return errt.NewIOError(err)
	}

	if err == io.ErrShortBuffer {
		return errt.NewIOError(err)
	}

	if err == io.ErrShortWrite {
		return errt.NewIOError(err)
	}

	if err == io.ErrUnexpectedEOF {
		return errt.NewIOError(err)
	}

	return NewMongoError(err)
}

// shouldTriggerRefresh checks whether an error should trigger the refreshing
// of the MongoDB connection, and whether the query that returned this error
// should be relaunched.
func shouldTriggerRefresh(err error) bool {
	return !(err == nil || mgo.IsDup(err) || err == mgo.ErrNotFound)
}
