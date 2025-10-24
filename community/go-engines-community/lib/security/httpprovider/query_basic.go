package httpprovider

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

// queryBasicProvider implements a Query Authentication provider.
// It validates user using security provider.
type queryBasicProvider struct {
	provider security.Provider
}

// NewQueryBasicProvider creates new provider.
// Deprecated : use Basic Auth.
func NewQueryBasicProvider(p security.Provider) security.HttpProvider {
	return &queryBasicProvider{provider: p}
}

func (p *queryBasicProvider) Auth(r *http.Request) (*security.User, error, bool) {
	q := r.URL.Query()
	username := q.Get("username")
	password := q.Get("password")
	if username == "" {
		return nil, nil, false
	}

	u, err := p.provider.Auth(r.Context(), username, password)

	return u, err, true
}
