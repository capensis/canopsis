package httpprovider

import (
	"errors"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/rs/zerolog"
)

// cookieProvider implements a Cookie Token Authentication provider.
// It must be used only for file access.
type cookieProvider struct {
	tokenProvider security.TokenProvider

	cookieName string
	logger     zerolog.Logger
}

func NewCookieProvider(
	tokenProvider security.TokenProvider,
	cookieName string,
	logger zerolog.Logger,
) security.HttpProvider {
	return &cookieProvider{
		tokenProvider: tokenProvider,

		cookieName: cookieName,
		logger:     logger,
	}
}

func (p *cookieProvider) Auth(r *http.Request) (*security.User, error, bool) {
	cookie, err := r.Cookie(p.cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil, false
		}
		return nil, err, false
	}

	tokenString, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err, false
	}

	user, err := p.tokenProvider.Auth(r.Context(), tokenString)
	return user, err, true
}
