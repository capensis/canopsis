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

func GetResponseItem(ar *fastjson.Arena, id string, status int, rawUser, err *fastjson.Value) *fastjson.Value {
	item := ar.NewObject()
	item.Set("status", ar.NewNumberInt(status))
	item.Set("item", rawUser)

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
