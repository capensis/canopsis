package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/ajg/form"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	amqplib "github.com/streadway/amqp"
	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type API interface {
	Send(c *gin.Context)
}

type api struct {
	publisher                   amqp.Publisher
	alarmCollection             mongo.DbCollection
	isAllowChangeSeverityToInfo bool
	logger                      zerolog.Logger
}

func NewApi(
	publisher amqp.Publisher,
	client mongo.DbClient,
	isAllowChangeSeverityToInfo bool,
	logger zerolog.Logger,
) API {
	return &api{
		publisher:                   publisher,
		isAllowChangeSeverityToInfo: isAllowChangeSeverityToInfo,
		alarmCollection:             client.Collection(mongo.AlarmMongoCollection),
		logger:                      logger,
	}
}

// Event structure used with swagger
type Event struct {
	Connector     string `json:"connector" example:"test_connector"`
	ConnectorName string `json:"connector_name" example:"test_connectorname"`
	SourceType    string `json:"source_type" example:"resource"`
	EventType     string `json:"event_type" example:"check"`
	Component     string `json:"component,omitempty" example:"test_component"`
	State         string `json:"state,omitempty" example:"1"`
	Resource      string `json:"resource" example:"test_resource"`
}

// Response structure used with swagger
type Response struct {
	SentEvents []Event `json:"sent_events"`
	// FailedEvents is an empty array left for compatibility with old handler
	FailedEvents []interface{} `json:"failed_events"`
	// RetryEvents is an empty array left for compatibility with old handler
	RetryEvents []interface{} `json:"retry_events"`
}

// Send event/events
// @Summary Send event/events
// @Description Send event/events
// @Tags events
// @ID event-send
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body Event true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ErrorResponse
// @Router /event [post]
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
		err = fmt.Errorf("the body should be an object or an array")
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
		eventType == types.EventTypeKeepstate ||
		eventType == types.EventTypeStatStateInterval ||
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

		if (eventType == types.EventTypeChangestate || eventType == types.EventTypeKeepstate) && state == 0 {
			if !api.isAllowChangeSeverityToInfo {
				api.logger.Err(fmt.Errorf("cannot set state to info with changestate/keepstate")).Str("event", string(value.MarshalTo(nil))).Msg("Event API error")
				return false
			}
		}
	}

	if eventType == types.EventTypeAck ||
		eventType == types.EventTypeAckremove ||
		eventType == types.EventTypeCancel ||
		eventType == types.EventTypeComment ||
		eventType == types.EventTypeUncancel ||
		eventType == types.EventTypeDeclareTicket ||
		eventType == types.EventTypeDone ||
		eventType == types.EventTypeAssocTicket ||
		eventType == types.EventTypeChangestate ||
		eventType == types.EventTypeKeepstate ||
		eventType == types.EventTypeSnooze ||
		eventType == types.EventTypeStatusIncrease ||
		eventType == types.EventTypeStatusDecrease ||
		eventType == types.EventTypeStateIncrease ||
		eventType == types.EventTypeStateDecrease {

		role, ok := c.Get(auth.RoleKey)
		if !ok {
			role = ""
			api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg("Cannot retrieve role from user")
		}

		value.Set("role", fastjson.MustParse(fmt.Sprintf("%q", role.(string))))
		api.logger.Info().Str("event", string(value.MarshalTo(nil))).Msgf("Role added to the event. event_type = %s, role = %s", eventType, role)
	}

	longOutputValue := value.Get("long_output")
	if longOutputValue != nil && longOutputValue.Type() != fastjson.TypeString {
		value.Set("long_output", fastjson.MustParse(`""`))
		api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msgf("Long output field is not a string : %s. Replacing it by \"\"", longOutputValue.Type())
	}

	author, err := getStringField(value, "author")
	if err != nil && !errors.Is(err, ErrFieldNotExists) {
		api.logger.Warn().Str("event", string(value.MarshalTo(nil))).Msg(err.Error())
		return false
	}

	if author == "" {
		userID := c.MustGet(auth.UserKey).(string)
		value.Set("author", fastjson.MustParse(fmt.Sprintf("%q", userID)))
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

	var alarm types.Alarm
	err = api.alarmCollection.FindOne(c.Request.Context(), bson.M{"d": eid}).Decode(&alarm)
	if err != nil && err != mongodriver.ErrNoDocuments {
		api.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to get alarm from mongo")
		return false
	}

	if err != mongodriver.ErrNoDocuments {
		processArray(value, "ma_parents", alarm.Value.Parents)
		processArray(value, "ma_children", alarm.Value.Children)

		if alarm.IsMetaAlarm() {
			cursor, err := api.alarmCollection.Aggregate(
				c.Request.Context(),
				[]bson.M{
					{
						"$match": bson.M{
							"d": bson.M{
								"$in": alarm.Value.Children,
							},
						},
					},
					{
						"$unwind": "$v.parents",
					},
					{
						"$group": bson.M{
							"_id": 1,
							"related_parents": bson.M{
								"$addToSet": bson.M{
									"$cond": bson.M{
										"if": bson.M{"$ne": bson.A{"$v.parents", alarm.EntityID}},
										"then": "$v.parents",
										"else": "$$REMOVE",
									},
								},
							},
						},
					},
				},
			)
			if err != nil {
				api.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to get related parents info from mongo")
				return false
			}
			defer cursor.Close(c.Request.Context())

			var relatedParentsInfo struct{
				RelatedParents []string `bson:"related_parents"`
			}

			if cursor.Next(c.Request.Context()) {
				err = cursor.Decode(&relatedParentsInfo)
				if err != nil {
					api.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to get related parents info from mongo")
					return false
				}

				processArray(value, "ma_related_parents", relatedParentsInfo.RelatedParents)
			}
		}
	}

	err = api.publisher.Publish(
		canopsis.CanopsisEventsExchange,
		"",
		false,
		false,
		amqplib.Publishing{
			ContentType:  "application/json",
			Body:         value.MarshalTo(nil),
			DeliveryMode: amqplib.Persistent,
		},
	)
	if err != nil {
		api.logger.Err(err).Str("event", string(value.MarshalTo(nil))).Msg("Failed to publish event")
		return false
	}

	return true
}

func processArray(value *fastjson.Value, key string, values []string) {
	if value.Exists(key) {
		return
	}
	items := fastjson.MustParse(`[]`)
	for idx, item := range values {
		items.SetArrayItem(idx, fastjson.MustParse(fmt.Sprintf("%q", item)))
	}
	value.Set(key, items)
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
