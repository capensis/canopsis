package middleware

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// SessionAuth returns a Session Authorization middleware.
// It checks session and retrieves user using provider.
// It checks auth only if session exists.
// Deprecated : don't use session.
func SessionAuth(db mongo.DbClient, configProvider config.ApiConfigProvider, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, security.SessionKey)

		if err != nil {
			if errors.As(err, &securecookie.MultiError{}) ||
				errors.Is(err, libsession.ErrNoSession) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
				return
			}
			panic(err)
		}

		if val, ok := session.Values["user"]; ok {
			if userID, ok := val.(string); ok {
				provider := userprovider.NewMongoProvider(db, configProvider)
				user, err := provider.FindByID(c, userID)

				if err != nil {
					panic(err)
				}

				if user == nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
					return
				}

				// The user credentials was found, set user's id to key UserKey in this context,
				// the user's id can be read later using c.MustGet(auth.UserKey).
				c.Set(auth.UserKey, user.ID)
				c.Set(auth.Username, user.DisplayName)
				c.Set(auth.RolesKey, user.Roles)
				c.Set(auth.ApiKey, user.AuthApiKey)
			} else {
				panic("user key is not string")
			}
		}

		c.Next()
	}
}
