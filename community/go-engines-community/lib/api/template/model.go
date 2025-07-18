package template

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"

type Response struct {
	IsValid bool                 `json:"is_valid"`
	Err     *validator.ErrReport `json:"err,omitempty"`
}
