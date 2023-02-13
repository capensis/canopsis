package pbehaviortype

import "errors"

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

var ErrLinkedTypeToPbehavior = ValidationError{Err: errors.New("type is linked to pbehavior")}
var ErrLinkedTypeToException = ValidationError{Err: errors.New("type is linked to exception")}
var ErrLinkedToActionType = ValidationError{Err: errors.New("type is linked to action")}
var ErrDefaultType = ValidationError{Err: errors.New("type is default")}
var ErrorDuplicatePriority = ValidationError{Err: errors.New("duplicate priority value")}
