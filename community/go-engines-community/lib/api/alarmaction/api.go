package alarmaction

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	BulkAck(c *gin.Context)
	BulkAckRemove(c *gin.Context)
	BulkSnooze(c *gin.Context)
	BulkCancel(c *gin.Context)
	BulkUncancel(c *gin.Context)
	BulkAssocTicket(c *gin.Context)
	BulkComment(c *gin.Context)
	BulkChangeState(c *gin.Context)
	AddBookmark(c *gin.Context)
	RemoveBookmark(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
	logger         zerolog.Logger
}

func NewApi(store Store, errorResponder httperror.Responder, logger zerolog.Logger) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

// Ack
// @Param body body AckRequest true "body"
func (a *api) Ack(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	request := AckRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Ack(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := Request{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.AckRemove(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := SnoozeRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Snooze(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := Request{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Cancel(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := Request{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Uncancel(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := AssocTicketRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.AssocTicket(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := CommentRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Comment(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := ChangeStateRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.ChangeState(c, c.Param("id"), request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// BulkAck
// @Param body body []BulkAckRequestItem true "body"
func (a *api) BulkAck(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkAckRequestItem) (string, error) {
		ok, err := a.store.Ack(c, request.ID, request.AckRequest, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkAckRemove
// @Param body body []BulkRequestItem true "body"
func (a *api) BulkAckRemove(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkRequestItem) (string, error) {
		ok, err := a.store.AckRemove(c, request.ID, request.Request, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkSnooze
// @Param body body []BulkSnoozeRequestItem true "body"
func (a *api) BulkSnooze(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkSnoozeRequestItem) (string, error) {
		ok, err := a.store.Snooze(c, request.ID, request.SnoozeRequest, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkCancel
// @Param body body []BulkRequestItem true "body"
func (a *api) BulkCancel(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkRequestItem) (string, error) {
		ok, err := a.store.Cancel(c, request.ID, request.Request, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkUncancel
// @Param body body []BulkRequestItem true "body"
func (a *api) BulkUncancel(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkRequestItem) (string, error) {
		ok, err := a.store.Uncancel(c, request.ID, request.Request, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkAssocTicket
// @Param body body []BulkAssocTicketRequestItem true "body"
func (a *api) BulkAssocTicket(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkAssocTicketRequestItem) (string, error) {
		ok, err := a.store.AssocTicket(c, request.ID, request.AssocTicketRequest, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkComment
// @Param body body []BulkCommentRequestItem true "body"
func (a *api) BulkComment(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkCommentRequestItem) (string, error) {
		ok, err := a.store.Comment(c, request.ID, request.CommentRequest, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

// BulkChangeState
// @Param body body []BulkChangeStateRequestItem true "body"
func (a *api) BulkChangeState(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	bulk.Handler(c, func(request BulkChangeStateRequestItem) (string, error) {
		ok, err := a.store.ChangeState(c, request.ID, request.ChangeStateRequest, userID, username)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

func (a *api) AddBookmark(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	found, err := a.store.AddBookmark(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) RemoveBookmark(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	found, err := a.store.RemoveBookmark(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
