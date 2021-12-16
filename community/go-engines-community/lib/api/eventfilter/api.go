package eventfilter

import (
	"encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"github.com/gin-gonic/gin/binding"
	"github.com/valyala/fastjson"
	"net/http"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type api struct {
	store        Store
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) common.BulkCrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// Create eventfilter
// @Summary Create eventfilter
// @Description Create eventfilter
// @Tags eventfilters
// @ID eventfilters-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EventFilter true "body"
// @Success 201 {object} EventFilter
// @Failure 400 {object} common.ErrorResponse
// @Router /eventfilter/rules [post]
func (a api) Create(c *gin.Context) {
	request := EventFilter{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.store.Insert(c.Request.Context(), &request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeEventFilter,
		ValueID:   request.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, request)
}

// Find all eventfilter
// @Summary Find all eventfilter
// @Description Get paginated list of eventfilter
// @Tags eventfilters
// @ID eventfilters-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Success 200 {object} common.PaginatedListResponse{data=[]EventFilter}
// @Failure 400 {object} common.ErrorResponse
// @Router /eventfilter/rules [get]
func (a api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), query)
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

// Get eventfilter by id
// @Summary Get eventfilter by id
// @Description Get eventfilter by id
// @Tags eventfilters
// @ID eventfilters-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "eventfilter id"
// @Success 200 {object} EventFilter
// @Failure 404 {object} common.ErrorResponse
// @Router /eventfilter/rules/{id} [get]
func (a api) Get(c *gin.Context) {
	evf, err := a.store.GetById(c.Request.Context(), c.Param("id"))

	if err == mongodriver.ErrNoDocuments || evf == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, evf)
}

// Update eventfilter by id
// @Summary Update eventfilter by id
// @Description Update eventfilter by id
// @Tags eventfilters
// @ID eventfilters-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "eventfilter id"
// @Param body body EventFilter true "body"
// @Success 200 {object} EventFilter
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /eventfilter/rules/{id} [put]
func (a api) Update(c *gin.Context) {
	var request EventFilterPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	var data EventFilter
	data.EventFilterPayload = request
	data.ID = c.Param("id")
	ok, _ := a.store.Update(c.Request.Context(), &data)

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err := a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEventFilter,
		ValueID:   data.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, data)
}

// Delete eventfilter by id
// @Summary Delete eventfilter by id
// @Description Delete eventfilter by id
// @Tags eventfilters
// @ID eventfilters-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "eventfilter id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /eventfilter/rules/{id} [delete]
func (a api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeEventFilter,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
}

func (a *api) BulkCreate(c *gin.Context) {
	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawUsers, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()
	logEntries := make([]logger.LogEntry, 0, len(rawUsers))

	for idx, rawUser := range rawUsers {
		userObject, err := rawUser.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		var eventFilterRequest EventFilter
		err = json.Unmarshal(userObject.MarshalTo(nil), &eventFilterRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(eventFilterRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, common.NewValidationErrorFastJsonValue(&ar, err, eventFilterRequest)))
			continue
		}

		err = a.store.Insert(c.Request.Context(), &eventFilterRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, eventFilterRequest.ID, http.StatusOK, rawUser, nil))
		logEntries = append(logEntries, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeEventFilter,
			ValueID:   eventFilterRequest.ID,
		})
	}

	err = a.actionLogger.BulkAction(c, logEntries)
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

func (a *api) BulkUpdate(c *gin.Context) {
	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawUsers, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()
	logEntries := make([]logger.LogEntry, 0, len(rawUsers))

	for idx, rawUser := range rawUsers {
		userObject, err := rawUser.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		var eventFilterRequest BulkUpdateRequestItem
		err = json.Unmarshal(userObject.MarshalTo(nil), &eventFilterRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(eventFilterRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, common.NewValidationErrorFastJsonValue(&ar, err, eventFilterRequest)))
			continue
		}

		eventFilter := EventFilter(eventFilterRequest)

		ok, err := a.store.Update(c.Request.Context(), &eventFilter)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawUser, ar.NewString("Not found")))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, eventFilter.ID, http.StatusOK, rawUser, nil))
		logEntries = append(logEntries, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeEventFilter,
			ValueID:   eventFilter.ID,
		})
	}

	err = a.actionLogger.BulkAction(c, logEntries)
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

func (a *api) BulkDelete(c *gin.Context) {
	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawUsers, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()
	logEntries := make([]logger.LogEntry, 0, len(rawUsers))

	for idx, rawUser := range rawUsers {
		userObject, err := rawUser.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		var ruleRequest BulkDeleteRequestItem
		err = json.Unmarshal(userObject.MarshalTo(nil), &ruleRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(ruleRequest)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, common.NewValidationErrorFastJsonValue(&ar, err, ruleRequest)))
			continue
		}

		ok, err := a.store.Delete(c.Request.Context(), ruleRequest.ID)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawUser, ar.NewString(err.Error())))
			continue
		}

		if !ok {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawUser, ar.NewString("Not found")))
			continue
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, ruleRequest.ID, http.StatusOK, rawUser, nil))
		logEntries = append(logEntries, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeEventFilter,
			ValueID:   ruleRequest.ID,
		})
	}

	err = a.actionLogger.BulkAction(c, logEntries)
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}
