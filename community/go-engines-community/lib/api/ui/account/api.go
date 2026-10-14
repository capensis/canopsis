package account

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

type API interface {
	Me(c *gin.Context)
	Update(c *gin.Context)
}

func NewApi(store Store, errorResponder response.ErrorResponder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

type api struct {
	store          Store
	errorResponder response.ErrorResponder
}

// Me
// @Success 200 {object} User
func (a *api) Me(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	user, err := a.store.GetOneBy(c, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if user == nil {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}

	c.JSON(http.StatusOK, user)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} User
func (a *api) Update(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := EditRequest{
		ID: userID,
		// author is needed for action logs, in that case the user modifies himself, so he's the author.
		Author: userID,
	}

	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	user, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if user == nil {
		a.errorResponder.Respond(c, httperror.ErrUnauthorized)

		return
	}

	c.JSON(http.StatusOK, user)
}
