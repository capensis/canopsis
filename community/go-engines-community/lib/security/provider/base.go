// provider contains authentication methods.
package provider

import (
	"context"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
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

func (p *baseProvider) GetName() string {
	return ""
}

func (p *baseProvider) Auth(ctx context.Context, username, password string) (*security.User, error) {
	user, err := p.userProvider.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err)
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
