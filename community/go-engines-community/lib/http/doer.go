package http

//go:generate mockgen -destination=../../mocks/lib/http/doer.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http Doer

import "net/http"

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
