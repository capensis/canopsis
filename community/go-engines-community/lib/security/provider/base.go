// provider contains authentication methods.
package provider

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/security"
	"git.canopsis.net/canopsis/go-engines/lib/security/password"
	"strings"
)

// baseProvider implements password-based authentication.
type baseProvider struct {
	userProvider    security.UserProvider
	passwordEncoder password.Encoder
}

// NewBaseProvider creates new provider.
func NewBaseProvider(p security.UserProvider, e password.Encoder) security.Provider {
	return &baseProvider{
		userProvider:    p,
		passwordEncoder: e,
	}
}

func (p *baseProvider) Auth(username, password string) (*security.User, error) {
	user, err := p.userProvider.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %v", err)
	}

	if user == nil || !user.IsEnabled {
		return nil, nil
	}

	hashedPassword := strings.ToLower(user.HashedPassword)
	if !p.passwordEncoder.IsValidPassword([]byte(hashedPassword), []byte(password)) {
		return nil, nil
	}

	return user, nil
}
