package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

const headerAuthorization = "Authorization"
const headerCookie = "Cookie"
const headerSetCookie = "set-cookie"

// ReverseProxyHandler directs requests to old API.
// It doesn't support old API session and uses auth api key for authentication.
func ReverseProxyHandler(legacyURL *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = legacyURL.Scheme
				req.URL.Host = legacyURL.Host
			},
			Transport: noCookieTransport{},
		}

		setProxyAuth(c)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

type noCookieTransport struct{}

// RoundTrip executes a single HTTP transaction and deletes cookies from response.
func (noCookieTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	// Ignore old API cookie
	if err == nil {
		resp.Header.Del(headerSetCookie)
	}

	return resp, err
}

// setProxyAuth deletes auth credentials of original request
// and sets auth api key from context to request.
func setProxyAuth(c *gin.Context) {
	// Delete all auth credentials
	c.Request.Header.Del(headerAuthorization)
	c.Request.Header.Del(headerCookie)
	c.Request.Header.Del(libsecurity.HeaderApiKey)
	query := c.Request.URL.Query()
	query.Del(libsecurity.QueryParamApiKey)
	c.Request.URL.RawQuery = query.Encode()

	// Add proxy auth credentials
	apiKey, ok := c.Get(auth.ApiKey)
	if ok {
		if s, ok := apiKey.(string); ok {
			c.Request.Header.Add(libsecurity.HeaderApiKey, s)
		}
	}
}
