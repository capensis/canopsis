package entitybasic

import "errors"

var ErrLinkedEntityToAlarm = errors.New("entity is linked to alarm")
var ErrComponent = errors.New("component has resources")
