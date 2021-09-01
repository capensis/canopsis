package httpprovider

import (
	"encoding/base64"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"

// basicProvider implements a Basic HTTP Authentication provider.
// It validates user using security provider.
type basicProvider struct {
	provider security.Provider
}

// NewBasicProvider creates new provider.
func NewBasicProvider(p security.Provider) security.HttpProvider {
	return &basicProvider{provider: p}
}

func (p *basicProvider) Auth(r *http.Request) (*security.User, error, bool) {
	header := r.Header.Get(headerAuthorization)
	if header == "" {
		return nil, nil, false
	}

	username, password := decodeHeader(header)
	if username == "" {
		return nil, nil, true
	}

	u, err := p.provider.Auth(r.Context(), username, password)

	return u, err, true
}

// decodeHeader retrieves username and password from header.
func decodeHeader(header string) (string, string) {
	header = strings.ReplaceAll(header, "Basic ", "")
	base, err := base64.StdEncoding.DecodeString(header)

	if err != nil {
		panic(err)
	}

	pair := strings.Split(string(base), ":")

	if len(pair) != 2 {
		return "", ""
	}

	return pair[0], pair[1]
}
