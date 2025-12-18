package pbehaviorexception

import (
	"errors"
	"net/http"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type API interface {
	crud.API
	Import(c *gin.Context)
	BulkDelete(c *gin.Context)
	BulkHide(c *gin.Context)
	BulkUnhide(c *gin.Context)
}

func NewApi(
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	maxFileSize uint64,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:          store,
		computeChan:    computeChan,
		maxFileSize:    maxFileSize,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

type api struct {
	store          Store
	computeChan    chan<- rpc.PbehaviorRecomputeEvent
	maxFileSize    uint64
	errorResponder httperror.Responder
	logger         zerolog.Logger
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusCreated, res)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	isLinked, err := a.store.IsLinked(c, res.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if isLinked {
		a.sendComputeTask()
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	exception, err := a.store.GetByID(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if exception == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, exception)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, c.Param("id"), userID)

	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Import
// @Success 200 {object} Response
func (a *api) Import(c *gin.Context) {
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			err = validation.NewSingleError("required", "file", "file", nil)
		}

		a.errorResponder.Respond(c, err)

		return
	}

	defer f.Close()

	name := c.Request.FormValue("name")
	pbhType := c.Request.FormValue("type")
	valErrs := make(validator.ValidationErrors, 0)
	if a.maxFileSize > 0 && uint64(fh.Size) > a.maxFileSize {
		valErrs = append(valErrs, validation.NewFieldErrorWithParam("filesize", "file", "file", strconv.FormatUint(a.maxFileSize, 10)))
	}

	if name == "" {
		valErrs = append(valErrs, validation.NewFieldError("required", "name", "name"))
	}

	if pbhType == "" {
		valErrs = append(valErrs, validation.NewFieldError("required", "type", "type"))
	}

	if len(valErrs) > 0 {
		a.errorResponder.Respond(c, validation.NewError(valErrs, nil))

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	exception, err := a.store.Import(c, name, pbhType, userID, f, fh)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, exception)
}

func (a *api) sendComputeTask() {
	a.computeChan <- rpc.PbehaviorRecomputeEvent{}
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, request.Author)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		return request.ID, nil
	}, a.errorResponder)
}

// BulkHide
// @Param body body []BulkToggleHiddenRequestItem true "body"
func (a *api) BulkHide(c *gin.Context) {
	a.toggleHidden(c, true)
}

// BulkUnhide
// @Param body body []BulkToggleHiddenRequestItem true "body"
func (a *api) BulkUnhide(c *gin.Context) {
	a.toggleHidden(c, false)
}

func (a *api) toggleHidden(c *gin.Context, hidden bool) {
	bulk.Handler(c, func(request BulkToggleHiddenRequestItem) (string, error) {
		found, err := a.store.ToggleHidden(c, request, hidden)
		if err != nil {
			return "", err
		}

		if !found {
			return "", httperror.ErrNotFound
		}

		return request.ID, nil
	}, a.errorResponder)
}
