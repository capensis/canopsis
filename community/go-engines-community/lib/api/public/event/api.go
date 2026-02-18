package event

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	libresponse "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin/binding"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

type API struct {
	publisher libamqp.Publisher
	responder libresponse.Responder
	logger    zerolog.Logger
}

func NewAPI(
	publisher libamqp.Publisher,
	responder libresponse.Responder,
	logger zerolog.Logger,
) *API {
	return &API{
		publisher: publisher,
		responder: responder,
		logger:    logger,
	}
}

func (a *API) Send(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var jsonValue *fastjson.Value
	if mediatype, _, err := mime.ParseMediaType(r.Header.Get("content-type")); err == nil && mediatype == binding.MIMEPOSTForm {
		if err = r.ParseForm(); err != nil {
			a.responder.RespondError(w, r, validation.NewInvalidRequestBodyError(err))

			return
		}

		var ar fastjson.Arena
		jsonValue = ar.NewObject()
		for k, v := range r.Form {
			if len(v) == 1 {
				jsonValue.Set(k, ar.NewString(v[0]))
			} else {
				av := ar.NewArray()
				for i := range v {
					av.SetArrayItem(i, ar.NewString(v[i]))
				}

				jsonValue.Set(k, av)
			}
		}
	} else if r.Body != nil {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			a.responder.RespondError(w, r, err)

			return
		}

		jsonValue, err = fastjson.ParseBytes(raw)
		if err != nil {
			a.responder.RespondError(w, r, validation.NewInvalidRequestBodyError(err))

			return
		}
	} else {
		a.responder.RespondError(w, r, validation.NewInvalidRequestBodyError(errors.New("empty request")))

		return
	}

	response := fastjson.MustParse(`{}`)
	sendEvents := fastjson.MustParse(`[]`)
	failedEvents := fastjson.MustParse(`[]`)
	retryEvents := fastjson.MustParse(`[]`)

	response.Set("sent_events", sendEvents)
	response.Set("failed_events", failedEvents)
	response.Set("retry_events", retryEvents)

	contextUser, err := authctx.GetUserKey(ctx)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	contextAuthor, err := authctx.GetUsername(ctx)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	roles, err := authctx.GetRoles(ctx)
	if err != nil {
		a.responder.RespondError(w, r, err)

		return
	}

	switch jsonValue.Type() {
	case fastjson.TypeObject:
		if !a.processValue(ctx, jsonValue, contextUser, contextAuthor, roles) {
			failedEvents.SetArrayItem(0, jsonValue)
			break
		}

		sendEvents.SetArrayItem(0, jsonValue)
	case fastjson.TypeArray:
		values, err := jsonValue.Array()
		if err != nil {
			a.responder.RespondError(w, r, validation.NewInvalidRequestBodyError(err))

			return
		}

		var sentIdx, failedIdx int

		for _, value := range values {
			if !a.processValue(ctx, value, contextUser, contextAuthor, roles) {
				failedEvents.SetArrayItem(failedIdx, value)
				failedIdx++

				continue
			}

			sendEvents.SetArrayItem(sentIdx, value)
			sentIdx++
		}
	default:
		a.responder.RespondError(w, r, validation.NewInvalidRequestBodyError(errors.New("unsupported json value")))

		return
	}

	w.Header().Set("Content-Type", libresponse.MIMEJSON)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response.MarshalTo(nil))
	if err != nil {
		a.logger.Err(err).Msg("failed to write response")
	}
}

func (a *API) processValue(ctx context.Context, value *fastjson.Value, contextUser, contextAuthor string, roles []string) bool {
	eventType, err := getStringField(value, "event_type")
	if err != nil {
		a.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	switch eventType {
	case types.EventTypeCheck, types.EventTypeMetaAlarm, types.EventTypeChangestate,
		types.EventTypeJunitTestSuiteUpdated:
		state, isNotInt, err := getIntField(value, "state")

		if err != nil {
			a.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
			return false
		}

		if isNotInt {
			var ar fastjson.Arena
			value.Set("state", ar.NewNumberInt(state))
		}
	}

	if eventType == types.EventTypeAck ||
		eventType == types.EventTypeAckremove ||
		eventType == types.EventTypeCancel ||
		eventType == types.EventTypeComment ||
		eventType == types.EventTypeUncancel ||
		eventType == types.EventTypeAssocTicket ||
		eventType == types.EventTypeChangestate ||
		eventType == types.EventTypeSnooze {

		if len(roles) > 0 {
			role := roles[0]
			value.Set("role", fastjson.MustParse(fmt.Sprintf("%q", role)))
			a.logger.Info().Str("event", string(value.MarshalTo(nil))).Msgf("Role added to the event. event_type = %s, role = %s", eventType, role)
		}
	}

	longOutputValue := value.Get("long_output")
	if longOutputValue != nil && longOutputValue.Type() != fastjson.TypeString {
		value.Set("long_output", fastjson.MustParse(`""`))
		a.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msgf("Long output field is not a string : %s. Replacing it by \"\"", longOutputValue.Type())
	}

	author, err := getStringField(value, "author")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		a.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	if author == "" {
		value.Set("author", fastjson.MustParse(fmt.Sprintf("%q", contextAuthor)))
	}

	user, err := getStringField(value, "user_id")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		a.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	if user == "" {
		value.Set("user_id", fastjson.MustParse(fmt.Sprintf("%q", contextUser)))
	}

	sourceType := types.SourceTypeConnector
	connector, err := getStringField(value, "connector")
	if err != nil || connector == "" {
		a.logger.Warn().Err(err).Str("key", "connector").Msg("")
		return false
	}

	connectorName, err := getStringField(value, "connector_name")
	if err != nil || connectorName == "" {
		a.logger.Warn().Err(err).Str("key", "connector_name").Msg("")
		return false
	}

	component, err := getStringField(value, "component")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		a.logger.Warn().Err(err).Str("key", "component").Msg("")
		return false
	}

	resource, err := getStringField(value, "resource")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		a.logger.Warn().Err(err).Str("key", "resource").Msg("")
		return false
	}

	if component == "" {
		if resource != "" {
			a.logger.Warn().Str("key", "component").Msg("resource is defined but component is empty")
		}
	} else if resource == "" {
		sourceType = types.SourceTypeComponent
	} else {
		sourceType = types.SourceTypeResource
	}

	if sourceType == types.SourceTypeConnector && (eventType == types.EventTypeCheck ||
		eventType == types.EventTypeContextUpdate) {
		a.logger.Warn().Str("key", "source_type").Msg("cannot create check event for connector")
		return false
	}

	eventSourceType, err := getStringField(value, "source_type")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		a.logger.Warn().Err(err).Str("key", "source_type").Msg("")
		return false
	}

	if eventSourceType != sourceType {
		value.Set("source_type", fastjson.MustParse(fmt.Sprintf("%q", sourceType)))
		a.logger.Info().
			Str("event", string(value.MarshalTo(nil))).
			Str("from", eventSourceType).
			Str("to", sourceType).
			Msgf("SourceType changed in the event")
	}

	err = a.publisher.PublishWithContext(
		ctx,
		canopsis.EventsExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         value.MarshalTo(nil),
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		a.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to publish event")
		return false
	}

	return true
}

func getStringField(value *fastjson.Value, key string) (string, error) {
	fieldValue := value.Get(key)
	if fieldValue == nil {
		return "", ErrFieldNotExists
	}

	if fieldValue.Type() != fastjson.TypeString {
		return "", ErrFieldWrongType
	}

	return string(fieldValue.GetStringBytes()), nil
}

func getIntField(value *fastjson.Value, key string) (int, bool, error) {
	fieldValue := value.Get(key)
	if fieldValue == nil {
		return 0, false, ErrFieldNotExists
	}

	if fType := fieldValue.Type(); fType != fastjson.TypeNumber {
		// try to convert string to int
		if fType == fastjson.TypeString {
			v, err := strconv.Atoi(string(fieldValue.GetStringBytes()))
			return v, true, err
		}
		return 0, true, ErrFieldWrongType
	}

	return fieldValue.GetInt(), false, nil
}
