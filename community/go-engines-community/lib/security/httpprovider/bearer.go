package httpprovider

import (
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

const bearerPrefix = "Bearer"

// bearerProvider implements a Bearer Token Authentication provider.
type bearerProvider struct {
	providers []security.TokenProvider
}

// NewBearerProvider creates new provider.
func NewBearerProvider(
	providers []security.TokenProvider,
) security.HttpProvider {
	return &bearerProvider{
		providers: providers,
	}
}

func (p *bearerProvider) Auth(r *http.Request) (*security.User, error, bool) {
	header := r.Header.Get(headerAuthorization)
	if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return nil, nil, false
	}

	tokenString := strings.TrimSpace(header[len(bearerPrefix):])
	for _, provider := range p.providers {
		user, err := provider.Auth(r.Context(), tokenString)
		if err != nil || user != nil {
			return user, err, true
		}
	}

	return nil, nil, true
}
