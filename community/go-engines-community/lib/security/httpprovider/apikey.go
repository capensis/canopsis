// httpprovider contains http authentication methods.
package httpprovider

import (
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

// apikeyProvider implements a ApiKey HTTP Authentication provider.
// It uses apiKey from request query or from header and retrieves user using user provider.
type apikeyProvider struct {
	userProvider security.UserProvider
}

// NewApikeyProvider creates new provider.
// Deprecated : use JWT token instead.
func NewApikeyProvider(p security.UserProvider) security.HttpProvider {
	return &apikeyProvider{userProvider: p}
}

func (p *apikeyProvider) Auth(r *http.Request) (*security.User, error, bool) {
	apiKey := r.URL.Query().Get(security.QueryParamApiKey)
	if apiKey == "" {
		apiKey = r.Header.Get(security.HeaderApiKey)
	}

	if apiKey == "" {
		return nil, nil, false
	}

	user, err := p.userProvider.FindByAuthApiKey(r.Context(), apiKey)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err), true
	}

	if user == nil || !user.IsEnabled {
		return nil, nil, true
	}

	return user, nil, true
}
