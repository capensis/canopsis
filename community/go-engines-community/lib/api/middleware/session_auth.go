package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/security"
	libsession "git.canopsis.net/canopsis/go-engines/lib/security/session"
	"git.canopsis.net/canopsis/go-engines/lib/security/userprovider"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// SessionAuth returns a Session Authorization middleware.
// It checks session and retrieves user using provider.
// It checks auth only if session exists.
func SessionAuth(db mongo.DbClient, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, security.SessionKey)

		if err != nil {
			_, isSessionError := err.(securecookie.MultiError)
			if isSessionError || err == libsession.ErrNoSession {
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
				return
			}
			panic(err)
		}

		if val, ok := session.Values["user"]; ok {
			if userId, ok := val.(string); ok {
				provider := userprovider.NewMongoProvider(db)
				user, err := provider.FindByID(userId)

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
				c.Set(auth.ApiKey, user.AuthApiKey)
			} else {
				panic("user key is not string")
			}
		}

		c.Next()
	}
}
