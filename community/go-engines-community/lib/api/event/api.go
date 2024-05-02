package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/ajg/form"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type API interface {
	Send(c *gin.Context)
}

type api struct {
	publisher       libamqp.Publisher
	alarmCollection mongo.DbCollection
	logger          zerolog.Logger
}

func NewApi(
	publisher libamqp.Publisher,
	client mongo.DbClient,
	logger zerolog.Logger,
) API {
	return &api{
		publisher:       publisher,
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
		logger:          logger,
	}
}

func (api *api) Send(c *gin.Context) {
	var err error
	var raw []byte
	var values []*fastjson.Value

	raw, err = c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	if mediatype, _, err := mime.ParseMediaType(c.GetHeader("content-type")); err == nil && mediatype == binding.MIMEPOSTForm {
		var u map[string]interface{}
		d := form.NewDecoder(bytes.NewBuffer(raw))
		if err := d.Decode(&u); err != nil {
			panic(err)
		}

		raw, err = json.Marshal(u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}
	response := fastjson.MustParse(`{}`)
	sendEvents := fastjson.MustParse(`[]`)
	failedEvents := fastjson.MustParse(`[]`)
	retryEvents := fastjson.MustParse(`[]`)

	response.Set("sent_events", sendEvents)
	response.Set("failed_events", failedEvents)
	response.Set("retry_events", retryEvents)

	switch jsonValue.Type() {
	case fastjson.TypeObject:
		if !api.processValue(c, jsonValue) {
			failedEvents.SetArrayItem(0, jsonValue)
			break
		}

		sendEvents.SetArrayItem(0, jsonValue)
	case fastjson.TypeArray:
		values, err = jsonValue.Array()
		if err != nil {
			break
		}

		var sentIdx, failedIdx int

		for _, value := range values {
			if !api.processValue(c, value) {
				failedEvents.SetArrayItem(failedIdx, value)
				failedIdx++

				continue
			}

			sendEvents.SetArrayItem(sentIdx, value)
			sentIdx++
		}
	default:
		err = errors.New("the body should be an object or an array")
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, gin.MIMEJSON, response.MarshalTo(nil))
}

func (api *api) processValue(c *gin.Context, value *fastjson.Value) bool {
	eventType, err := getStringField(value, "event_type")
	if err != nil {
		api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	if eventType == types.EventTypeCheck ||
		eventType == types.EventTypeMetaAlarm ||
		eventType == types.EventTypeChangestate ||
		eventType == types.EventTypeJunitTestSuiteUpdated {
		state, isNotInt, err := getIntField(value, "state")

		if err != nil {
			api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
			return false
		}

		if isNotInt {
			var a fastjson.Arena
			value.Set("state", a.NewNumberInt(state))
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

		roles, ok := c.Get(auth.RolesKey)
		role := ""
		if ok {
			if s, ok := roles.([]string); ok && len(s) > 0 {
				role = s[0]
			}
		} else {
			api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg("Cannot retrieve role from user")
		}

		value.Set("role", fastjson.MustParse(fmt.Sprintf("%q", role)))
		api.logger.Info().Str("event", string(value.MarshalTo(nil))).Msgf("Role added to the event. event_type = %s, role = %s", eventType, role)
	}

	longOutputValue := value.Get("long_output")
	if longOutputValue != nil && longOutputValue.Type() != fastjson.TypeString {
		value.Set("long_output", fastjson.MustParse(`""`))
		api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msgf("Long output field is not a string : %s. Replacing it by \"\"", longOutputValue.Type())
	}

	contextAuthor := c.MustGet(auth.Username).(string)
	contextUser := c.MustGet(auth.UserKey).(string)

	author, err := getStringField(value, "author")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	if author == "" {
		value.Set("author", fastjson.MustParse(fmt.Sprintf("%q", contextAuthor)))
	}

	user, err := getStringField(value, "user_id")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	if user == "" {
		value.Set("user_id", fastjson.MustParse(fmt.Sprintf("%q", contextUser)))
	}

	var eid string
	refRkValue, err := getStringField(value, "ref_rk")
	if err == nil {
		eid = refRkValue
	}

	if eid == "" {
		sourceType := types.SourceTypeConnector
		connector, err := getStringField(value, "connector")
		if err != nil || connector == "" {
			api.logger.Warn().Err(err).Str("key", "connector").Msg("")
			return false
		}
		connectorName, err := getStringField(value, "connector_name")
		if err != nil || connectorName == "" {
			api.logger.Warn().Err(err).Str("key", "connector_name").Msg("")
			return false
		}

		eid = fmt.Sprintf("%s/%s", connector, connectorName)

		component, err := getStringField(value, "component")
		if err != nil && !errors.Is(err, ErrFieldNotExists) {
			api.logger.Warn().Err(err).Str("key", "component").Msg("")
			return false
		}

		resource, err := getStringField(value, "resource")
		if err != nil && !errors.Is(err, ErrFieldNotExists) {
			api.logger.Warn().Err(err).Str("key", "resource").Msg("")
			return false
		}

		if component == "" {
			if resource != "" {
				api.logger.Warn().Str("key", "component").Msg("resource is defined but component is empty")
			}
		} else {
			if resource == "" {
				eid = component
				sourceType = types.SourceTypeComponent
			} else {
				eid = fmt.Sprintf("%s/%s", resource, component)
				sourceType = types.SourceTypeResource
			}
		}

		if sourceType == types.SourceTypeConnector && eventType == types.EventTypeCheck {
			api.logger.Warn().Str("key", "source_type").Msg("cannot create check event for connector")
			return false
		}

		eventSourceType, err := getStringField(value, "source_type")
		if err != nil && !errors.Is(err, ErrFieldNotExists) {
			api.logger.Warn().Err(err).Str("key", "source_type").Msg("")
			return false
		}

		if eventSourceType != sourceType {
			value.Set("source_type", fastjson.MustParse(fmt.Sprintf("%q", sourceType)))
			api.logger.Info().
				Str("event", string(value.MarshalTo(nil))).
				Str("from", eventSourceType).
				Str("to", sourceType).
				Msgf("SourceType changed in the event")
		}
	}

	err = api.alarmCollection.FindOne(c, bson.M{"d": eid}).Err()
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		api.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to get alarm from mongo")
		return false
	}

	err = api.publisher.PublishWithContext(
		c,
		canopsis.CanopsisEventsExchange,
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
		api.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to publish event")
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
