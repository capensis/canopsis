package logger

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"math"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"
	ActionApproval = "approval"
	ActionExport   = "export"
	ActionImport   = "import"
)

const ApprovalDecisionApprove = "approve"
const ApprovalDecisionDismiss = "dismiss"

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
	ValueTypeJobConfig          = "jobconfig"
	ValueTypeJob                = "job"
	ValueTypeInstruction        = "instruction"
	ValueTypeStateSetting       = "statesetting"
	ValueTypeBroadcastMessage   = "broadcastmessage"
	ValueAssociativeTable       = "associativetable"
	ValueTypeIdleRule           = "idlerule"

	ValueTypeResolveRule  = "resolverule"
	ValueTypeFlappingRule = "flappingrule"

	ValueTypeUserPreferences = "userpreferences"

	ValueTypeFilter = "filter"
)

type ActionLogger interface {
	Action(c *gin.Context, logEntry LogEntry) error
	BulkAction(c *gin.Context, logEntries []LogEntry) error
	Err(err error, msg string)
}

type logger struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	zLog         zerolog.Logger
}

type LogEntry struct {
	Action         string    `bson:"action"`
	ValueType      string    `bson:"value_type"`
	ValueID        string    `bson:"value_id"`
	Author         string    `bson:"author"`
	Time           time.Time `bson:"time"`
	ValidationTime time.Time `bson:"validation_time,omitempty"`
	ApproverUser   string    `bson:"approver_user,omitempty"`
	ApproverRole   string    `bson:"approver_role,omitempty"`
	Decision       string    `bson:"decision,omitempty"`
}

func NewActionLogger(dbClient mongo.DbClient, zLog zerolog.Logger) ActionLogger {
	return &logger{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.ActionLogMongoCollection),
		zLog:         zLog,
	}
}

func (l *logger) Action(c *gin.Context, logEntry LogEntry) error {
	if logEntry.Author == "" {
		userID := c.MustGet(auth.UserKey)
		logEntry.Author = userID.(string)
	}

	logEntry.Time = time.Now()

	l.zLog.Info().
		Str("action", logEntry.Action).
		Str("value_type", logEntry.ValueType).
		Str("value_id", logEntry.ValueID).
		Str("author", logEntry.Author).
		Str("time", logEntry.Time.String()).
		Msg("ActionLog: ")

	_, err := l.dbCollection.UpdateOne(c.Request.Context(), bson.M{"value_type": logEntry.ValueType, "value_id": logEntry.ValueID}, bson.M{"$set": logEntry}, options.Update().SetUpsert(true))
	return err
}

func (l *logger) BulkAction(c *gin.Context, logEntries []LogEntry) error {
	if len(logEntries) == 0 {
		return nil
	}

	var err error
	userID := c.MustGet(auth.UserKey).(string)
	ctx := c.Request.Context()

	writeModels := make([]mongodriver.WriteModel, 0, int(math.Min(float64(canopsis.DefaultBulkSize), float64(len(logEntries)))))
	now := time.Now()

	for _, e := range logEntries {
		e.Time = now
		if e.Author == "" {
			e.Author = userID
		}

		writeModels = append(
			writeModels,
			mongodriver.NewInsertOneModel().SetDocument(e),
		)

		l.zLog.Info().
			Str("action", e.Action).
			Str("value_type", e.ValueType).
			Str("value_id", e.ValueID).
			Str("author", e.Author).
			Str("time", e.Time.String()).
			Msg("ActionLog: ")

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = l.dbCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = l.dbCollection.BulkWrite(ctx, writeModels)
	}

	return err
}

func (l *logger) Err(err error, msg string) {
	l.zLog.Err(err).Msg(msg)
}
