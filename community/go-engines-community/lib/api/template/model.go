package template

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	IsValid bool              `json:"is_valid"`
	Report  *validator.Report `json:"report,omitempty"`
}
