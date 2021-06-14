package errt_test

import (
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/errt"
)

type myError struct {
	errt.ErrT
	moreInfos bool
}

type MyError interface {
	errt.ErrT
	IsMyError()
	MoreInfos() bool
}

// IsMyError function does nothing because it's only here to implement the MyError interface.
// This makes type assertion possible.
func (e myError) IsMyError() {
}

func (e myError) MoreInfos() bool {
	return e.moreInfos
}

func NewMyError(err error, moreInfos bool) MyError {
	// Always return nil if the root error is nil
	if err == nil {
		return nil
	}
	return myError{
		ErrT:      errt.NewErrT(err),
		moreInfos: moreInfos,
	}
}

func IReturnAnError() error {
	oerr := errors.New("this is my error")
	return NewMyError(oerr, true)
}

func Example() {
	err := IReturnAnError()

	if err != nil {
		switch referr := err.(type) {
		case MyError:
			fmt.Printf("i have more infos: %v", referr.MoreInfos())

		default:
			fmt.Printf("unknown error: %v", referr)
		}
	}
}
