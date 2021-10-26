package httpprovider

import (
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

const bearerPrefix = "Bearer"

// bearerProvider implements a Bearer Token Authentication provider.
type bearerProvider struct {
	provider security.TokenProvider
}

// NewBearerProvider creates new provider.
func NewBearerProvider(
	provider security.TokenProvider,
) security.HttpProvider {
	return &bearerProvider{
		provider: provider,
	}
}

func (p *bearerProvider) Auth(r *http.Request) (*security.User, error, bool) {
	header := r.Header.Get(headerAuthorization)
	if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return nil, nil, false
	}

	tokenString := strings.TrimSpace(header[len(bearerPrefix):])
	user, err := p.provider.Auth(r.Context(), tokenString)
	if err != nil {
		return nil, err, true
	}

	return user, nil, true
}
