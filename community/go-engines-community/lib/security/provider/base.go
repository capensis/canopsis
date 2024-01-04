// Package provider contains authentication methods.
package provider

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
)

// baseProvider implements password-based authentication.
type baseProvider struct {
	userProvider               security.UserProvider
	passwordEncoder            password.Encoder
	deprecatedPasswordEncoders []password.Encoder
}

// NewBaseProvider creates new provider.
func NewBaseProvider(
	p security.UserProvider,
	passwordEncoder password.Encoder,
	deprecatedPasswordEncoders ...password.Encoder,
) security.Provider {
	return &baseProvider{
		userProvider:               p,
		passwordEncoder:            passwordEncoder,
		deprecatedPasswordEncoders: deprecatedPasswordEncoders,
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

	bytesHashedPwd := []byte(user.HashedPassword)
	bytesPwd := []byte(password)

	if ok, _ := p.passwordEncoder.IsValidPassword(bytesHashedPwd, bytesPwd); ok {
		return user, nil
	}

	for _, passwordEncoder := range p.deprecatedPasswordEncoders {
		if ok, _ := passwordEncoder.IsValidPassword(bytesHashedPwd, bytesPwd); !ok {
			continue
		}

		newHash, err := p.passwordEncoder.EncodePassword(bytesPwd)
		if err != nil {
			return nil, err
		}

		user.HashedPassword = string(newHash)
		err = p.userProvider.UpdateHashedPassword(ctx, user.ID, user.HashedPassword)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}
