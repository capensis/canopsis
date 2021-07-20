package websocket

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// AuthorizeRoom determines if current subject has been authorized to subscribe to a room.
func AuthorizeRoom(enforcer security.Enforcer, roomPerms map[string][]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj := c.MustGet(auth.UserKey)
		room := c.Param("room")
		perms := roomPerms[room]
		if len(perms) == 0 {
			return
		}

		vals := []interface{}{subj}
		for _, v := range perms {
			vals = append(vals, v)
		}
		ok, err := enforcer.Enforce(vals...)

		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}

		c.Next()
	}
}
