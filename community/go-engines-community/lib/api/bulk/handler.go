package bulk

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

var ErrUnauthorized = errors.New("unauthorized")

type BadRequestError struct {
	Err error
}

func (v BadRequestError) Error() string {
	return v.Err.Error()
}

func Handler[T any](
	c *gin.Context,
	f func(T) (string, error),
	logger zerolog.Logger,
) {
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

	var arena fastjson.Arena
	response := arena.NewArray()
	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, arena.NewString(err.Error())))
			continue
		}

		var request T
		err = json.Unmarshal(object.MarshalTo(nil), &request)
		if err != nil {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, arena.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&arena, err, request)))
			continue
		}

		id, err := f(request)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusForbidden, rawObject, arena.NewString(common.ForbiddenResponse.Error)))
				continue
			}

			valErr := common.ValidationError{}
			if errors.As(err, &valErr) {
				response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&arena, valErr, request)))
				continue
			}

			badRequestError := BadRequestError{}
			if errors.As(err, &badRequestError) {
				response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, arena.NewString(err.Error())))
				continue
			}

			logger.Err(err).Msg("cannot process bulk item")
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusInternalServerError, rawObject, arena.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if id == "" {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusNotFound, rawObject, arena.NewString(common.NotFoundResponse.Error)))
			continue
		}

		response.SetArrayItem(idx, GetResponseItem(&arena, id, http.StatusOK, rawObject, nil))
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

func HandlerWithGrouping[T any](
	c *gin.Context,
	compareWithPrev func(prev, cur T) bool,
	mergeReqs func(merged, cur T) T,
	procReq func(T) (id string, err error),
	logger zerolog.Logger,
) {
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

	var arena fastjson.Arena
	response := arena.NewArray()
	type groupedRequest struct {
		Req        T
		RawObjects map[int]*fastjson.Value
	}

	groupedReqs := make([]groupedRequest, 0)
	var prevReq T
	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, arena.NewString(err.Error())))
			continue
		}

		var req T
		err = json.Unmarshal(object.MarshalTo(nil), &req)
		if err != nil {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, arena.NewString(err.Error())))
			continue
		}

		err = binding.Validator.ValidateStruct(req)
		if err != nil {
			response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&arena, err, req)))
			continue
		}

		if idx > 0 && compareWithPrev(prevReq, req) {
			groupedReqs[len(groupedReqs)-1].Req = mergeReqs(groupedReqs[len(groupedReqs)-1].Req, req)
			groupedReqs[len(groupedReqs)-1].RawObjects[idx] = rawObject
		} else {
			groupedReqs = append(groupedReqs, groupedRequest{Req: req, RawObjects: map[int]*fastjson.Value{idx: rawObject}})
		}

		prevReq = req
	}

	for _, g := range groupedReqs {
		req := g.Req
		id, err := procReq(req)
		if err != nil {
			for idx, rawObject := range g.RawObjects {
				if errors.Is(err, ErrUnauthorized) {
					response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusForbidden, rawObject, arena.NewString(common.ForbiddenResponse.Error)))
					continue
				}

				valErr := common.ValidationError{}
				if errors.As(err, &valErr) {
					response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&arena, valErr, req)))
					continue
				}

				logger.Err(err).Msg("cannot process bulk item")
				response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusInternalServerError, rawObject, arena.NewString(common.InternalServerErrorResponse.Error)))
			}

			continue
		}

		if id == "" {
			for idx, rawObject := range g.RawObjects {
				response.SetArrayItem(idx, GetResponseItem(&arena, "", http.StatusNotFound, rawObject, arena.NewString(common.NotFoundResponse.Error)))
			}

			continue
		}

		for idx, rawObject := range g.RawObjects {
			response.SetArrayItem(idx, GetResponseItem(&arena, id, http.StatusOK, rawObject, nil))
		}
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
}

func GetResponseItem(ar *fastjson.Arena, id string, status int, rawItem, err *fastjson.Value) *fastjson.Value {
	item := ar.NewObject()
	item.Set("status", ar.NewNumberInt(status))
	item.Set("item", rawItem)

	if err == nil {
		item.Set("id", ar.NewString(id))
		return item
	}

	if err.Type() == fastjson.TypeString {
		item.Set("error", err)
	}

	if err.Type() == fastjson.TypeObject {
		item.Set("errors", err)
	}

	return item
}
