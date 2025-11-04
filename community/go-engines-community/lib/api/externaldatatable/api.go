package externaldatatable

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	crud.API
	BulkDeleteData(c *gin.Context)
	GetSchema(c *gin.Context)

	Import(*gin.Context)
	ImportStatus(*gin.Context)
	ImportData(*gin.Context)
	Preview(c *gin.Context)
	ImportComplete(*gin.Context)

	Export(c *gin.Context)
	ExportStatus(c *gin.Context)
	ExportDownload(c *gin.Context)

	ListData(c *gin.Context)
	GetData(c *gin.Context)
	CreateData(c *gin.Context)
	UpdateData(c *gin.Context)
	DeleteData(c *gin.Context)
}

func NewAPI(
	store Store,
	importWorker ImportWorker,
	maxFileSize uint64,
	exportTaskCreator export.TaskCreator,
	exportParamsEncoder encoding.Encoder,
	logger zerolog.Logger,
) API {
	return &api{
		store:               store,
		importWorker:        importWorker,
		maxFileSize:         maxFileSize,
		exportTaskCreator:   exportTaskCreator,
		exportParamsEncoder: exportParamsEncoder,
		logger:              logger,
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
	}
}

type api struct {
	store               Store
	importWorker        ImportWorker
	maxFileSize         uint64
	exportTaskCreator   export.TaskCreator
	exportParamsEncoder encoding.Encoder
	logger              zerolog.Logger
	exportSeparators    map[string]rune
}

// Create
// @Param body body EditRequest true "body"
// @Success 200 {array} Table
func (a *api) Create(c *gin.Context) {
	r := EditRequest{}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, err := a.store.Create(c, r)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	c.JSON(http.StatusCreated, res)
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Table}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	aggregationResult, err := a.store.Find(c, r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {array} Table
func (a *api) Get(c *gin.Context) {
	res, err := a.store.FindOne(c, c.Param("table"))
	if err != nil {
		panic(err)
	}

	if res.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, res)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {array} Table
func (a *api) Update(c *gin.Context) {
	r := UpdateRequest{
		ID: c.Param("table"),
	}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, err := a.store.Update(c, r)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	if res.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("table"), c.MustGet(auth.UserKey).(string))
	if err != nil {
		if errors.Is(err, ErrConfigNotDeletable) || errors.Is(err, ErrLinkedNotDeletable) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

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

// Import
// @Success 200 {array} ImportJob
func (a *api) Import(c *gin.Context) {
	id := c.Param("table")
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationError("file", "File is missing.").ValidationErrorResponse())

			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: "request has invalid structure"})

		return
	}

	defer f.Close()
	valErrors := make(map[string]string)
	if a.maxFileSize > 0 && uint64(fh.Size) > a.maxFileSize {
		valErrors["file"] = fmt.Sprintf("File size %d exceeds limit %d", fh.Size, a.maxFileSize)
	}

	separatorStr := c.Request.FormValue("separator")
	separator, ok := a.exportSeparators[separatorStr]
	if separatorStr != "" && !ok {
		valErrors["separator"] = "Separator must be one of [comma semicolon tab space] or empty."
	}

	if len(valErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrors(valErrors).ValidationErrorResponse())

		return
	}

	job, err := a.importWorker.CreateImportJob(c, id, separator, f)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, job)
}

// Preview
// @Success 200 {array} ImportJob
func (a *api) Preview(c *gin.Context) {
	id := c.Param("id")

	var r PreviewRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	job, err := a.importWorker.CreatePreviewJob(c, id, r)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	if job.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, job)
}

// ImportStatus
// @Success 200 {array} ImportJob
func (a *api) ImportStatus(c *gin.Context) {
	job, err := a.importWorker.GetJob(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if job.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, job)
}

func (a *api) ImportData(c *gin.Context) {
	var r ListPreviewRequest
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	job, err := a.importWorker.GetJob(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if job.ID == "" || job.Status != ImportStatusSucceeded && job.Status != ImportStatusPreviewSucceeded && job.Status != ImportStatusPreviewFailed {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	aggregationResult, err := a.store.FindPreviewData(c, job, r)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

// ImportComplete
// @Param body body ImportCompleteRequest true "body"
func (a *api) ImportComplete(c *gin.Context) {
	r := ImportCompleteRequest{}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	ok, err := a.importWorker.CompleteJob(c, c.Param("id"), r.ColumnTags)
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

// Export
// @Param request body ExportRequest true "request"
// @Success 200 {object} ExportResponse
func (a *api) Export(c *gin.Context) {
	var r ExportRequest
	r.ID = c.Param("table")
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	t, err := a.store.FindOne(c, r.ID)
	if err != nil {
		panic(err)
	}

	if t.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	columns := t.getColumns()

	if len(r.SearchBy) > 0 || len(r.Fields) > 0 {
		hasCol := make(map[string]bool, len(columns))
		for _, column := range columns {
			hasCol[column] = true
		}

		valErrMsgs := make(map[string]string)
		for i, v := range r.SearchBy {
			if !hasCol[v] {
				valErrMsgs["search_by."+strconv.Itoa(i)] = "SearchBy"
			}
		}

		for i, f := range r.Fields {
			if !hasCol[f.Name] {
				valErrMsgs["fields."+strconv.Itoa(i)+".name"] = "Name"
			}
		}

		if len(valErrMsgs) > 0 {
			errMsg := " must be one of [" + strings.Join(columns, " ") + "]."
			for k := range valErrMsgs {
				valErrMsgs[k] += errMsg
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrors(valErrMsgs).ValidationErrorResponse())

			return
		}
	}

	fields := r.Fields
	if len(fields) == 0 {
		fields = make(export.Fields, len(columns))
		for i, v := range columns {
			fields[i] = export.Field{Name: v}
		}
	}

	separator := a.exportSeparators[r.Separator]
	params, err := a.exportParamsEncoder.Encode(r.ExportFetchParameters)
	if err != nil {
		panic(err)
	}

	userID := c.MustGet(auth.UserKey).(string)
	task, err := a.exportTaskCreator.Create(c, export.TaskParameters{
		Type:           "externaldata",
		Parameters:     string(params),
		Fields:         fields,
		Separator:      separator,
		FilenamePrefix: "externaldata",
		UserID:         userID,
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     task.ID,
		Status: task.Status,
	})
}

// ExportStatus
// @Success 200 {object} ExportResponse
func (a *api) ExportStatus(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportTaskCreator.Get(c, id)
	if err != nil {
		panic(err)
	}

	if t == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     id,
		Status: t.Status,
	})
}

func (a *api) ExportDownload(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportTaskCreator.Get(c, id)
	if err != nil {
		panic(err)
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/csv")
	c.FileAttachment(t.File, t.Filename)
}

func (a *api) CreateData(c *gin.Context) {
	r := make(map[string]any)
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, err := a.store.CreateData(c, c.Param("table"), r)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	if len(res) == 0 {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusCreated, res)
}

func (a *api) ListData(c *gin.Context) {
	var r ListDataRequest
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	table, err := a.store.FindOne(c, c.Param("table"))
	if err != nil {
		panic(err)
	}

	if table.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	aggregationResult, err := a.store.FindData(c, table.getDBTableName(), table.Type, table.ColumnConfigs, r)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) GetData(c *gin.Context) {
	res, err := a.store.FindOneData(c, c.Param("table"), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if len(res) == 0 {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) UpdateData(c *gin.Context) {
	r := make(map[string]any)
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, err := a.store.UpdateData(c, c.Param("table"), c.Param("id"), r)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	if len(res) == 0 {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) DeleteData(c *gin.Context) {
	table, err := a.store.FindOne(c, c.Param("table"))
	if err != nil {
		panic(err)
	}

	if table.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	ok, err := a.store.DeleteData(c, table, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	c.Status(http.StatusNoContent)
}

// BulkDeleteData
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDeleteData(c *gin.Context) {
	table, err := a.store.FindOne(c, c.Param("table"))
	if err != nil {
		panic(err)
	}

	if table.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.DeleteData(c, table, request.ID)
		if err != nil || !ok {
			return "", err
		}

		return request.ID, nil
	}, a.logger)
}

func (a *api) GetSchema(c *gin.Context) {
	t, err := a.store.FindOne(c, c.Param("table"))
	if err != nil {
		panic(err)
	}

	if t.ID == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err = w.Write(t.getColumns())
	if err != nil {
		panic(err)
	}

	w.Flush()
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, t.Name+".csv"))
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
