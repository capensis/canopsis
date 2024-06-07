package icon

import (
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"github.com/gin-gonic/gin"
)

const (
	websocketMsgTypeCreate = iota
	websocketMsgTypeUpdate
	websocketMsgTypeDelete
)

type API interface {
	common.CrudAPI
	Patch(c *gin.Context)
}

func NewApi(
	store Store,
	websocketHub websocket.Hub,
	maxSize int64,
	mimeTypes []string,
) API {
	return &api{
		store:        store,
		websocketHub: websocketHub,
		maxSize:      maxSize,
		mimeTypes:    mimeTypes,
	}
}

type api struct {
	store        Store
	websocketHub websocket.Hub
	maxSize      int64
	mimeTypes    []string
}

type websocketMsg struct {
	ID   string `json:"_id"`
	Type int    `json:"type"`
}

// Create
// @Success 200 {array} Response
func (a *api) Create(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		Author: userID,
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	mimeType, valErr := a.validateFile(request.File)
	if valErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
		return
	}

	request.MimeType = mimeType
	res, err := a.store.Create(c, request)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	a.websocketHub.Send(websocket.RoomIcons, websocketMsg{
		ID:   res.ID,
		Type: websocketMsgTypeCreate,
	})
	c.JSON(http.StatusCreated, res)
}

func (a *api) Get(c *gin.Context) {
	res, err := a.store.Get(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, res)
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
	userID := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		ID:     c.Param("id"),
		Author: userID,
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	mimeType, valErr := a.validateFile(request.File)
	if valErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
		return
	}

	request.MimeType = mimeType
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

	a.websocketHub.Send(websocket.RoomIcons, websocketMsg{
		ID:   res.ID,
		Type: websocketMsgTypeUpdate,
	})
	c.JSON(http.StatusOK, res)
}

// Patch
// @Success 200 {object} Response
func (a *api) Patch(c *gin.Context) {
	userID := c.MustGet(auth.UserKey).(string)
	request := PatchRequest{
		ID:     c.Param("id"),
		Author: userID,
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	if request.File != nil {
		mimeType, valErr := a.validateFile(request.File)
		if valErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		request.MimeType = mimeType
	}

	res, err := a.store.Patch(c, request)
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

	a.websocketHub.Send(websocket.RoomIcons, websocketMsg{
		ID:   res.ID,
		Type: websocketMsgTypeUpdate,
	})
	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c, id, c.MustGet(auth.UserKey).(string))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.websocketHub.Send(websocket.RoomIcons, websocketMsg{
		ID:   id,
		Type: websocketMsgTypeDelete,
	})
	c.Status(http.StatusNoContent)
}

func (a *api) validateFile(file *multipart.FileHeader) (string, *common.ValidationError) {
	if file.Size > a.maxSize {
		err := common.NewValidationError("file", fmt.Sprintf("File size %d exceeds limit %d", file.Size, a.maxSize))
		return "", &err
	}

	mimeType := mime.TypeByExtension(path.Ext(file.Filename))
	if !slices.Contains(a.mimeTypes, mimeType) {
		err := common.NewValidationError("file", "Invalid mime type: "+mimeType)
		return "", &err
	}

	return mimeType, nil
}
