package action

import "errors"

var InvalidOldAlarmPattern = errors.New("old alarm pattern is invalid")
var InvalidOldEntityPattern = errors.New("old entity pattern is invalid")
var AlarmPatternError = errors.New("alarm pattern returned error")
var EntityPatternError = errors.New("alarm pattern returned error")
