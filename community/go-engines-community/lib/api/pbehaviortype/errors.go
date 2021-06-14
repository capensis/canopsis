package pbehaviortype

import "errors"

type ValidationError error

var ErrLinkedTypeToPbehavior = ValidationError(errors.New("type is linked to pbehavior"))
var ErrLinkedTypeToException = ValidationError(errors.New("type is linked to exception"))
var ErrLinkedToActionType = ValidationError(errors.New("type is linked to action"))
var ErrDefaultType = ValidationError(errors.New("type is default"))
var ErrorDuplicatePriority = ValidationError(errors.New("duplicate priority value"))
