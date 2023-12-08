package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// ReloadEnforcerPolicyOnChange loads security policy if request changes policy.
func ReloadEnforcerPolicyOnChange(enforcer security.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Writer.Status()
		switch s {
		case http.StatusOK, http.StatusNoContent, http.StatusCreated:
			err := enforcer.LoadPolicy()

			if err != nil {
				panic(err)
			}
		}

		c.Next()
	}
}
