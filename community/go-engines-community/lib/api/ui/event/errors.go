package event

import "errors"

var ErrFieldNotExists = errors.New("field doesn't exist")
var ErrFieldWrongType = errors.New("field has a wrong type")
