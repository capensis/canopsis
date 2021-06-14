package httpprovider

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"net/http"
)

// queryProvider implements a Query Authentication provider.
// It validates user using security provider.
type queryProvider struct {
	provider security.Provider
}

// NewQueryProvider creates new provider.
func NewQueryProvider(p security.Provider) security.HttpProvider {
	return &queryProvider{provider: p}
}

func (p *queryProvider) Auth(r *http.Request) (*security.User, error, bool) {
	q := r.URL.Query()
	username := q.Get("username")
	password := q.Get("password")
	if username == "" {
		return nil, nil, false
	}

	u, err := p.provider.Auth(username, password)

	return u, err, true
}
