package icon

import (
	"mime"
	"mime/multipart"
	"net/http"
	"path"
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	websocketMsgTypeCreate = iota
	websocketMsgTypeUpdate
	websocketMsgTypeDelete
)

type API interface {
	crud.API
	Patch(c *gin.Context)
}

func NewApi(
	store Store,
	websocketHub websocket.Hub,
	maxSize uint64,
	mimeTypes []string,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:          store,
		websocketHub:   websocketHub,
		maxSize:        maxSize,
		mimeTypes:      mimeTypes,
		errorResponder: errorResponder,
	}
}

type api struct {
	store          Store
	websocketHub   websocket.Hub
	maxSize        uint64
	mimeTypes      []string
	errorResponder httperror.Responder
}

type websocketMsg struct {
	ID   string `json:"_id"`
	Type int    `json:"type"`
}

// Create
// @Success 200 {array} Response
func (a *api) Create(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := EditRequest{
		Author: userID,
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	mimeType, err := a.validateFile(request.File)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	request.MimeType = mimeType
	res, err := a.store.Create(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
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
		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, res)
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	query := pagination.FilteredQuery{}
	query.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.List(c, query)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Update
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := EditRequest{
		ID:     c.Param("id"),
		Author: userID,
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	mimeType, err := a.validateFile(request.File)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	request.MimeType = mimeType
	res, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := PatchRequest{
		ID:     c.Param("id"),
		Author: userID,
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if request.File != nil {
		request.MimeType, err = a.validateFile(request.File)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}
	}

	res, err := a.store.Patch(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

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
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, id, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	a.websocketHub.Send(websocket.RoomIcons, websocketMsg{
		ID:   id,
		Type: websocketMsgTypeDelete,
	})
	c.Status(http.StatusNoContent)
}

func (a *api) validateFile(file *multipart.FileHeader) (string, error) {
	if uint64(file.Size) > a.maxSize {
		return "", validation.NewError(
			validator.ValidationErrors{validation.NewFieldErrorWithParam("filesize", "file", "file", strconv.FormatUint(a.maxSize, 10))},
			nil,
		)
	}

	mimeType := mime.TypeByExtension(path.Ext(file.Filename))
	if !slices.Contains(a.mimeTypes, mimeType) {
		return "", validation.NewError(
			validator.ValidationErrors{validation.NewFieldErrorWithParam("filetype", "file", "file", strings.Join(a.mimeTypes, " "))},
			nil,
		)
	}

	return mimeType, nil
}
