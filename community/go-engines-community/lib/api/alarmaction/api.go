package alarmaction

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	Ack(c *gin.Context)
	AckRemove(c *gin.Context)
	Snooze(c *gin.Context)
	Cancel(c *gin.Context)
	Uncancel(c *gin.Context)
	AssocTicket(c *gin.Context)
	Comment(c *gin.Context)
	ChangeState(c *gin.Context)
}

type api struct {
	store Store
}

func NewApi(store Store) API {
	return &api{store: store}
}

// Ack
// @Param body body AckRequest true "body"
func (a *api) Ack(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := AckRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Ack(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// AckRemove
// @Param body body Request true "body"
func (a *api) AckRemove(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := Request{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.AckRemove(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// Snooze
// @Param body body SnoozeRequest true "body"
func (a *api) Snooze(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := SnoozeRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Snooze(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// Cancel
// @Param body body Request true "body"
func (a *api) Cancel(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := Request{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Cancel(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// Uncancel
// @Param body body Request true "body"
func (a *api) Uncancel(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := Request{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Uncancel(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// AssocTicket
// @Param body body AssocTicketRequest true "body"
func (a *api) AssocTicket(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := AssocTicketRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.AssocTicket(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// Comment
// @Param body body CommentRequest true "body"
func (a *api) Comment(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := CommentRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.Comment(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangeState
// @Param body body ChangeStateRequest true "body"
func (a *api) ChangeState(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := ChangeStateRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.store.ChangeState(c, c.Param("id"), request, userID)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
