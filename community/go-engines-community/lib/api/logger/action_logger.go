package logger

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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
	Action(ctx context.Context, userID string, logEntry LogEntry) error
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

func (l *logger) Action(ctx context.Context, userID string, logEntry LogEntry) error {
	if logEntry.Author == "" {
		logEntry.Author = userID
	}

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
