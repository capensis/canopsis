package logger

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionExport = "export"
	ActionImport = "import"
)

const (
	ValueTypeUser               = "user"
	ValueTypeRole               = "role"
	ValueTypeView               = "view"
	ValueTypeViewGroup          = "viewgroup"
	ValueTypePlayList           = "playlist"
	ValueTypeEventFilter        = "eventfilter"
	ValueTypeScenario           = "scenario"
	ValueTypeMetaalarmRule      = "metaalarmrule"
	ValueTypeDynamicInfo        = "dynamicinfo"
	ValueTypeEntity             = "entity"
	ValueTypeEntityService      = "entityservice"
	ValueTypeEntityCategory     = "entitycategory"
	ValueTypePbehaviorType      = "pbehaviortype"
	ValueTypePbehaviorReason    = "pbehaviorreason"
	ValueTypePbehaviorException = "pbehaviorexception"
	ValueTypePbehavior          = "pbehavior"
	ValueTypeHeartbeat          = "heartbeat"
	ValueTypeJobConfig          = "jobconfig"
	ValueTypeJob                = "job"
	ValueTypeInstruction        = "instruction"
	ValueTypeStateSetting       = "statesetting"
	ValueTypeBroadcastMessage   = "broadcastmessage"
	ValueAssociativeTable       = "associativetable"
	ValueTypeIdleRule           = "idlerule"
)

type ActionLogger interface {
	Action(c *gin.Context, logEntry LogEntry) error
	Err(err error, msg string)
}

type logger struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	zLog         zerolog.Logger
}

type LogEntry struct {
	Action    string    `bson:"action"`
	ValueType string    `bson:"value_type"`
	ValueID   string    `bson:"value_id"`
	Author    string    `bson:"author"`
	Time      time.Time `bson:"time"`
}

func NewActionLogger(dbClient mongo.DbClient, zLog zerolog.Logger) ActionLogger {
	return &logger{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.ActionLogMongoCollection),
		zLog:         zLog,
	}
}

func (l *logger) Action(c *gin.Context, logEntry LogEntry) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userID := c.MustGet(auth.UserKey)

	logEntry.Author = userID.(string)
	logEntry.Time = time.Now()

	l.zLog.Info().
		Str("action", logEntry.Action).
		Str("value_type", logEntry.ValueType).
		Str("value_id", logEntry.ValueID).
		Str("author", logEntry.Author).
		Str("time", logEntry.Time.String()).
		Msg("ActionLog: ")

	_, err := l.dbCollection.UpdateOne(ctx, bson.M{"value_type": logEntry.ValueType, "value_id": logEntry.ValueID}, bson.M{"$set": logEntry}, options.Update().SetUpsert(true))
	return err
}

func (l *logger) Err(err error, msg string) {
	l.zLog.Err(err).Msg(msg)
}
