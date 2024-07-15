package oauth

const (
	oauthSessionPrefix = "oauth-session-"
)

type oidcClaims struct {
	ScopesSupported []string `json:"scopes_supported"`
	ClaimsSupported []string `json:"claims_supported"`
}

func (c *oidcClaims) ValidateScopes(scopes []string) bool {
	supported := make(map[string]bool)
	for _, v := range c.ScopesSupported {
		supported[v] = true
	}

	for _, v := range scopes {
		if !supported[v] {
			return false
		}
	}

	return true
}

type loginRequest struct {
	Redirect string `form:"redirect" binding:"required,url"`
}
