package icon

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file"
	"github.com/gin-gonic/gin"
)

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	maxSize int64,
	mimeTypes []string,
) common.CrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,
		maxSize:      maxSize,
		mimeTypes:    mimeTypes,
	}
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
	maxSize      int64
	mimeTypes    []string
}

// Create
// @Success 200 {array} Response
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	valErr := a.validateFile(&request)
	if valErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
		return
	}

	res, err := a.store.Create(c, request)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeIcon,
		ValueID:   res.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, res)
}

func (a *api) Get(c *gin.Context) {
	m, err := a.store.Get(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if m == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Header("Etag", fmt.Sprintf("%q", m.Etag))
	c.Header("Content-Type", m.MimeType)
	filename := ""
	exts, err := mime.ExtensionsByType(m.MimeType)
	if err == nil && len(exts) > 0 {
		filename = file.Sanitize(m.Title) + exts[0]
	}

	c.FileAttachment(a.store.GetFilepath(*m), filename)
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	query := pagination.FilteredQuery{}
	query.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.List(c, query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Update
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	valErr := a.validateFile(&request)
	if valErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
		return
	}

	res, err := a.store.Update(c, request)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeIcon,
		ValueID:   res.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) validateFile(r *EditRequest) *common.ValidationError {
	if r.File.Size > a.maxSize {
		err := common.NewValidationError("file", fmt.Sprintf("File size %d exceeds limit %d", r.File.Size, a.maxSize))
		return &err
	}

	r.MimeType = mime.TypeByExtension(path.Ext(r.File.Filename))
	if !slices.Contains(a.mimeTypes, r.MimeType) {
		err := common.NewValidationError("file", "Invalid mime type: "+r.MimeType)
		return &err
	}

	return nil
}
