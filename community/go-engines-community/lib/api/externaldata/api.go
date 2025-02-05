package externaldata

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	common.CrudAPI
	Import(*gin.Context)
	ImportStatus(*gin.Context)
	ImportData(*gin.Context)
	ImportComplete(*gin.Context)
	ListData(c *gin.Context)
	GetData(c *gin.Context)
	CreateData(c *gin.Context)
	UpdateData(c *gin.Context)
	DeleteData(c *gin.Context)
	BulkDeleteData(c *gin.Context)
	GetSchema(c *gin.Context)
}

func NewAPI(store Store, importWorker ImportWorker, maxFileSize uint64, logger zerolog.Logger) API {
	return &api{
		store:        store,
		importWorker: importWorker,
		maxFileSize:  maxFileSize,
		logger:       logger,
	}
}

type api struct {
	store        Store
	importWorker ImportWorker
	maxFileSize  uint64
	logger       zerolog.Logger
}

// Create
// @Param body body EditRequest true "body"
// @Success 200 {array} Response
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
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
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
// @Success 200 {array} Response
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
// @Param body body EditRequest true "body"
// @Success 200 {array} Response
func (a *api) Update(c *gin.Context) {
	r := EditRequest{
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
	id := c.Param("id")
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

	delimiterStr := c.Request.FormValue("delimiter")
	if delimiterStr == "" {
		valErrors["delimiter"] = "Delimiter is missing."
	} else if len(delimiterStr) > 1 {
		valErrors["delimiter"] = "Delimiter is too long."
	}

	if len(valErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrors(valErrors).ValidationErrorResponse())

		return
	}

	delimiter := rune(delimiterStr[0])
	job, err := a.importWorker.CreateJob(c, id, delimiter, f)
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
	var r ListDataRequest
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	job, err := a.importWorker.GetJob(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if job.ID == "" || job.Status != ImportStatusSucceeded {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	columns := make([]string, len(job.ColumnLengths))
	i := 0
	for col := range job.ColumnLengths {
		columns[i] = col
		i++
	}

	sort.Strings(columns)
	aggregationResult, err := a.store.FindData(c, job.Table, job.Type, columns, r)
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

	ok, err := a.importWorker.CompleteJob(c, c.Param("id"), r.ColumnTypes)
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

func (a *api) CreateData(c *gin.Context) {
	r := make(map[string]string)
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

	columns := make([]string, len(table.ColumnTypes))
	i := 0
	for col := range table.ColumnTypes {
		columns[i] = col
		i++
	}

	sort.Strings(columns)
	aggregationResult, err := a.store.FindData(c, table.Name, table.Type, columns, r)
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
	r := make(map[string]string)
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

	fields := make([]string, len(t.ColumnTypes))
	i, n := 0, len(t.ColumnTypes)
	for f := range t.ColumnTypes {
		fields[i] = f
		i++
		n += len(f)
	}

	sort.Strings(fields)
	b := &bytes.Buffer{}
	b.Grow(n)
	w := csv.NewWriter(b)
	err = w.Write(fields)
	if err != nil {
		panic(err)
	}

	w.Flush()
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, t.Name+".csv"))
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
