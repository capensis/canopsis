package externaldata

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

type API interface {
	common.CrudAPI
	Import(*gin.Context)
	ImportStatus(*gin.Context)
	ImportData(*gin.Context)
	ImportComplete(*gin.Context)
	ListData(c *gin.Context)
}

func NewAPI(store Store, importWorker ImportWorker) API {
	return &api{
		store:        store,
		importWorker: importWorker,
	}
}

type api struct {
	store        Store
	importWorker ImportWorker
}

// Create
// @Param body body CreateRequest true "body"
// @Success 200 {array} Response
func (a *api) Create(c *gin.Context) {
	request := CreateRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

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

	c.JSON(http.StatusCreated, res)
}

func (a *api) List(c *gin.Context) {
	var request ListRequest
	request.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	aggregationResult, err := a.store.Find(c, request)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(request.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Get(c *gin.Context) {
	panic("not implemented")
}

func (a *api) Update(c *gin.Context) {
	panic("not implemented")
}

func (a *api) Delete(c *gin.Context) {
	panic("not implemented")
}

// Import
// @Success 200 {array} ImportJob
func (a *api) Import(c *gin.Context) {
	id := c.Param("id")
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{
				"file": "File is missing.",
			}})

			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: "request has invalid structure"})

		return
	}

	defer f.Close()
	delimiterStr := c.Request.FormValue("delimiter")
	valErrors := make(map[string]string)
	if delimiterStr == "" {
		valErrors["delimiter"] = "Delimiter is missing."
	} else if len(delimiterStr) > 1 {
		valErrors["delimiter"] = "Delimiter is too long."
	}

	if len(valErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: valErrors})
		return
	}

	delimiter := rune(delimiterStr[0])
	job, err := a.importWorker.CreateJob(c, id, delimiter, f, fh)
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
	var request ListDataRequest
	request.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

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

	aggregationResult, err := a.store.FindData(c, job.Table, job.Type, job.Columns, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	res, err := common.NewPaginatedResponse(request.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

// ImportComplete
// @Param body body ImportCompleteRequest true "body"
func (a *api) ImportComplete(c *gin.Context) {
	request := ImportCompleteRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	ok, err := a.importWorker.CompleteJob(c, c.Param("id"), request.Columns)
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

func (a *api) ListData(c *gin.Context) {
	var request ListDataRequest
	request.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	table, err := a.store.FindOne(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if table.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	columns := make([]string, len(table.Columns))
	i := 0
	for col := range table.Columns {
		columns[i] = col
		i++
	}

	aggregationResult, err := a.store.FindData(c, table.Name, table.Type, columns, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())

			return
		}

		panic(err)
	}

	res, err := common.NewPaginatedResponse(request.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}
