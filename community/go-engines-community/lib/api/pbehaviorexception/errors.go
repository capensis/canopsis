package pbehaviorexception

import "errors"

var ErrTypeNotExists = errors.New("type doesn't exist")
var ErrLinkedException = errors.New("exception is linked with pbehavior")
