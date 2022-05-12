package alarm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

type API interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	GetDetails(c *gin.Context)
	ListManual(c *gin.Context)
	ListByService(c *gin.Context)
	Count(c *gin.Context)
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
}

type api struct {
	store                  Store
	exportExecutor         export.TaskExecutor
	defaultExportFields    export.Fields
	exportSeparators       map[string]rune
	timezoneConfigProvider config.TimezoneConfigProvider

	logger zerolog.Logger
}

func NewApi(
	store Store,
	executor export.TaskExecutor,
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
) API {
	fields := []string{"_id", "v.connector", "v.connector_name", "v.component",
		"v.resource", "v.output", "v.state.val", "v.status.val"}
	defaultExportFields := make(export.Fields, len(fields))
	for i, field := range fields {
		defaultExportFields[i] = export.Field{
			Name:  field,
			Label: field,
		}
	}

	return &api{
		store:               store,
		exportExecutor:      executor,
		defaultExportFields: defaultExportFields,
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
		timezoneConfigProvider: timezoneConfigProvider,
		logger:                 logger,
	}
}

// Find all alarms
// @Summary Find all alarms
// @Description Get paginated list of alarms. Use parameters "multi_sort[]=field_1,asc&multi_sort[]=field_2,desc&multi_sort[]=field_3,asc" to sort list by multiple fields. "sort_key", "sort_dir" are left for compatibility with older ways of sorting
// @Tags alarms
// @ID alarms-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ListRequestWithPagination true "request"
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /alarms [get]
func (a *api) List(c *gin.Context) {
	var r ListRequestWithPagination
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	apiKey := c.MustGet(auth.ApiKey).(string)

	aggregationResult, err := a.store.Find(c.Request.Context(), apiKey, r)
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

func (a *api) Get(c *gin.Context) {
	apiKey := c.MustGet(auth.ApiKey).(string)
	alarm, err := a.store.GetByID(c.Request.Context(), c.Param("id"), apiKey)
	if err != nil {
		panic(err)
	}

	if alarm == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, alarm)
}

func (a *api) GetDetails(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawObjects, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	apiKey := c.MustGet(auth.ApiKey).(string)
	defaultQuery := pagination.GetDefaultQuery()
	response := make([]DetailsResponse, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response[idx].Status = http.StatusBadRequest
			response[idx].Error = err.Error()
			continue
		}

		var request DetailsRequest
		err = json.Unmarshal(object.MarshalTo(nil), &request)
		if err != nil {
			response[idx].Status = http.StatusBadRequest
			response[idx].Error = err.Error()
			continue
		}

		request.Steps.Paginate = true
		if request.Steps.Limit == 0 {
			request.Steps.Limit = defaultQuery.Limit
		}

		request.Children.Paginate = true
		if request.Children.Limit == 0 {
			request.Children.Limit = defaultQuery.Limit
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response[idx].ID = request.ID
			response[idx].Status = http.StatusBadRequest
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				response[idx].Errors = common.TransformValidationErrors(errs, request).Errors
			} else {
				response[idx].Error = "Request has invalid structure"
			}
			continue
		}

		details, err := a.store.GetDetails(c.Request.Context(), apiKey, request)
		if err != nil {
			response[idx].ID = request.ID
			response[idx].Status = http.StatusInternalServerError
			response[idx].Error = common.InternalServerErrorResponse.Error
			a.logger.Err(err).Msg("cannot fetch alarm details")
			continue
		}

		if details == nil {
			response[idx].ID = request.ID
			response[idx].Status = http.StatusNotFound
			response[idx].Error = common.NotFoundResponse.Error
			continue
		}

		response[idx].ID = request.ID
		response[idx].Status = http.StatusOK
		response[idx].Data = *details
	}

	c.JSON(http.StatusMultiStatus, response)
}

func (a *api) ListManual(c *gin.Context) {
	var r ManualRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	alarms, err := a.store.FindManual(c.Request.Context(), r.Search)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, alarms)
}

func (a *api) ListByService(c *gin.Context) {
	r := pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	apiKey := c.MustGet(auth.ApiKey).(string)
	aggregationResult, err := a.store.FindByService(c.Request.Context(), c.Param("id"), apiKey, r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Count alarms
// @Summary Count alarms
// @Description Get counts of alarms
// @Tags alarms
// @ID alarms-get-counts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query FilterRequest true "request"
// @Success 200 {object} Count
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /alarm-counters [get]
func (a *api) Count(c *gin.Context) {
	var r FilterRequest

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	res, err := a.store.Count(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// Start export alarms
// @Summary Start export alarms
// @Description Start export alarms
// @Tags alarms
// @ID alarms-export-start
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request body ExportRequest true "request"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /alarm-export [post]
func (a *api) StartExport(c *gin.Context) {
	var r ExportRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	separator := a.exportSeparators[r.Separator]
	exportFields := r.Fields
	if len(exportFields) == 0 {
		exportFields = a.defaultExportFields
	}

	apiKey := c.MustGet(auth.ApiKey).(string)

	taskID, err := a.exportExecutor.StartExecute(c.Request.Context(), export.Task{
		Filename:     "alarms",
		ExportFields: exportFields,
		Separator:    separator,
		DataFetcher: getDataFetcher(a.store, apiKey, r, exportFields.Fields(),
			a.timezoneConfigProvider.Get().Location),
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     taskID,
		Status: export.TaskStatusRunning,
	})
}

// Get status of export alarms
// @Summary Get status of export alarms
// @Description Get status of export alarms
// @Tags alarms
// @ID alarms-export-get
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "export task id"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /alarm-export/{id} [get]
func (a *api) GetExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(c.Request.Context(), id)
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

// Get result of export alarms
// @Summary Get result of export alarms
// @Description Get result of export alarms
// @Tags alarms
// @ID alarms-export-download
// @Produce text/csv
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "export task id"
// @Success 200 {object} http.Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /alarm-export/{id}/download [get]
func (a *api) DownloadExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, t.Filename))
	c.Header("Content-Type", "text/csv")
	c.ContentType()
	c.File(t.File)
}
