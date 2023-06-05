package pbehaviorreason

import "errors"

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

var ErrLinkedReasonToPbehavior = ValidationError{Err: errors.New("reason is linked to pbehavior")}
var ErrLinkedReasonToAction = ValidationError{Err: errors.New("reason is linked to action")}
