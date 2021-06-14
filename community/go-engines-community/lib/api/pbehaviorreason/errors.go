package pbehaviorreason

import "errors"

type ValidationError error

var ErrLinkedReasonToPbehavior = ValidationError(errors.New("reason is linked to pbehavior"))
var ErrLinkedReasonToAction = ValidationError(errors.New("reason is linked to action"))
