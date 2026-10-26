package alarmtimeline

import (
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/gin-gonic/gin"
)

type API interface {
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
}

type api struct {
	store            Store
	taskCreator      export.TaskCreator
	exportSeparators map[string]rune
	encoder          encoding.Encoder
	errorResponder   response.ErrorResponder
}

func NewApi(
	store Store,
	taskCreator export.TaskCreator,
	encoder encoding.Encoder,
	errorResponder response.ErrorResponder,
) API {
	return &api{
		store:       store,
		taskCreator: taskCreator,
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
		encoder:        encoder,
		errorResponder: errorResponder,
	}
}

// StartExport
// @Param request body ExportRequest true "request"
// @Success 200 {object} ExportResponse
func (a *api) StartExport(c *gin.Context) {
	var r ExportRequest
	if err := binding.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	exists, err := a.store.Exists(c, r.AlarmID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !exists {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	params, err := a.encoder.Encode(r.ExportFetchParameters)
	if err != nil {
		a.errorResponder.Respond(c, fmt.Errorf("encode timeline export parameters: %w", err))

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	task, err := a.taskCreator.Create(c, export.TaskParameters{
		Type:       ExportType,
		Parameters: string(params),
		Fields: export.Fields{
			{Name: "_t", Label: "type"},
			{Name: "t", Label: "timestamp"},
			{Name: "a", Label: "author"},
			{Name: "m", Label: "message"},
			{Name: "val", Label: "value"},
			{Name: "initiator", Label: "initiator"},
		},
		Separator:      a.exportSeparators[r.Separator],
		TimeFormat:     r.TimeFormat,
		FilenamePrefix: ExportType,
		UserID:         userID,
	})
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     task.ID,
		Status: task.Status,
	})
}

// GetExport
// @Success 200 {object} ExportResponse
func (a *api) GetExport(c *gin.Context) {
	id := c.Param("id")

	t, err := a.taskCreator.Get(c, id)
	if err != nil {
		a.errorResponder.Respond(c, fmt.Errorf("fetch timeline export task %q: %w", id, err))

		return
	}

	if t == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     id,
		Status: t.Status,
	})
}

func (a *api) DownloadExport(c *gin.Context) {
	id := c.Param("id")

	t, err := a.taskCreator.Get(c, id)
	if err != nil {
		a.errorResponder.Respond(c, fmt.Errorf("fetch timeline export task %q: %w", id, err))

		return
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/csv")
	c.FileAttachment(t.File, t.Filename)
}
