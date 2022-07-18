package pbehavior

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
)

var ErrReasonNotExists = errors.New("reason doesn't exist")
var ErrExceptionNotExists = errors.New("exception doesn't exist")

type ValidationError struct {
	field string
	err   error
}

func (v ValidationError) Error() string {
	return v.err.Error()
}

func (v ValidationError) ToResponse() common.ValidationErrorResponse {
	return common.ValidationErrorResponse{
		Errors: map[string]string{
			v.field: v.Error(),
		},
	}
}
