package pattern

import "errors"

var ErrInvalidOldAlarmPattern = errors.New("old alarm pattern is invalid")
var ErrInvalidOldEntityPattern = errors.New("old entity pattern is invalid")
var ErrInvalidOldEventPattern = errors.New("old event pattern is invalid")
