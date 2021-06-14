package dev

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// ReloadEnforcerPolicy loads security policy on each request.
// It's required because fixtures are loaded before each functional test suit.
func ReloadEnforcerPolicy(enforcer security.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := enforcer.LoadPolicy()

		if err != nil {
			panic(err)
		}

		c.Next()
	}
}
