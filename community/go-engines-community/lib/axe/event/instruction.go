package event

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewInstructionProcessor(
	client mongo.DbClient,
	metricsSender metrics.Sender,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &instructionProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
		metricsSender:   metricsSender,
		amqpPublisher:   amqpPublisher,
		encoder:         encoder,
		logger:          logger,
		alarmStepTypeMap: map[string]string{
			// Manual instruction
			types.EventTypeInstructionStarted:   types.AlarmStepInstructionStart,
			types.EventTypeInstructionPaused:    types.AlarmStepInstructionPause,
			types.EventTypeInstructionResumed:   types.AlarmStepInstructionResume,
			types.EventTypeInstructionCompleted: types.AlarmStepInstructionComplete,
			types.EventTypeInstructionFailed:    types.AlarmStepInstructionFail,
			// Auto instruction
			types.EventTypeAutoInstructionStarted:   types.AlarmStepAutoInstructionStart,
			types.EventTypeAutoInstructionCompleted: types.AlarmStepAutoInstructionComplete,
			types.EventTypeAutoInstructionFailed:    types.AlarmStepAutoInstructionFail,
			// Manual and auto instruction
			types.EventTypeInstructionAborted: types.AlarmStepInstructionAbort,
		},
		alarmChangeTypeMap: map[string]types.AlarmChangeType{
			// Manual instruction
			types.EventTypeInstructionStarted:   types.AlarmChangeTypeInstructionStart,
			types.EventTypeInstructionPaused:    types.AlarmChangeTypeInstructionPause,
			types.EventTypeInstructionResumed:   types.AlarmChangeTypeInstructionResume,
			types.EventTypeInstructionCompleted: types.AlarmChangeTypeInstructionComplete,
			types.EventTypeInstructionFailed:    types.AlarmChangeTypeInstructionFail,
			// Auto instruction
			types.EventTypeAutoInstructionStarted:   types.AlarmChangeTypeAutoInstructionStart,
			types.EventTypeAutoInstructionCompleted: types.AlarmChangeTypeAutoInstructionComplete,
			types.EventTypeAutoInstructionFailed:    types.AlarmChangeTypeAutoInstructionFail,
			// Manual and auto instruction
			types.EventTypeInstructionAborted: types.AlarmChangeTypeInstructionAbort,
			// Job
			types.EventTypeInstructionJobStarted:   types.AlarmChangeTypeInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmChangeTypeInstructionJobComplete,
			types.EventTypeInstructionJobFailed:    types.AlarmChangeTypeInstructionJobFail,
		},
	}
}

type instructionProcessor struct {
	alarmCollection    mongo.DbCollection
	metricsSender      metrics.Sender
	amqpPublisher      libamqp.Publisher
	encoder            encoding.Encoder
	logger             zerolog.Logger
	alarmStepTypeMap   map[string]string
	alarmChangeTypeMap map[string]types.AlarmChangeType
}

func (p *instructionProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	alarmStepType := p.alarmStepTypeMap[event.EventType]
	alarmChangeType, ok := p.alarmChangeTypeMap[event.EventType]
	if !ok {
		return result, nil
	}

	var match bson.M
	set := bson.M{}
	if alarmStepType == "" {
		match = getOpenAlarmMatch(event)
	} else {
		match = getOpenAlarmMatchWithStepsLimit(event)
		newStepQuery := execStepUpdateQueryWithInPbhInterval(alarmStepType, "", event.Parameters.Output, event.Parameters)
		set["v.steps"] = addStepUpdateQuery(newStepQuery)
	}

	switch alarmChangeType {
	case types.AlarmStepAutoInstructionStart:
		set["v.inactive_start"] = bson.M{"$cond": bson.M{
			"if":   "$auto_instruction_in_progress",
			"then": event.Parameters.Timestamp,
			"else": "$v.inactive_start",
		}}
		set["v.inactive_duration"] = bson.M{"$sum": bson.A{
			"$v.inactive_duration",
			bson.M{"$cond": bson.M{
				"if": bson.M{"$and": []bson.M{
					{"$eq": bson.A{"$auto_instruction_in_progress", true}},
					{"$gt": bson.A{"$v.inactive_start", 0}},
				}},
				"then": bson.M{"$subtract": bson.A{
					event.Parameters.Timestamp,
					"$v.inactive_start",
				}},
				"else": 0,
			}},
		}}
	case types.AlarmStepInstructionComplete, types.AlarmStepInstructionFail:
		set["kpi_executed_instructions"] = bson.M{"$concatArrays": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$kpi_executed_instructions",
				"then": "$kpi_executed_instructions",
				"else": bson.A{},
			}},
			bson.M{"$cond": bson.M{
				"if": bson.M{"$and": []bson.M{
					{"$in": bson.A{
						event.Parameters.Instruction,
						bson.M{"$cond": bson.M{
							"if":   "$kpi_assigned_instructions",
							"then": "$kpi_assigned_instructions",
							"else": bson.A{},
						}},
					}},
					{"$not": bson.M{"$in": bson.A{
						event.Parameters.Instruction,
						bson.M{"$cond": bson.M{
							"if":   "$kpi_executed_instructions",
							"then": "$kpi_executed_instructions",
							"else": bson.A{},
						}},
					}}},
				}},
				"then": bson.A{event.Parameters.Instruction},
				"else": bson.A{},
			}},
		}}
	case types.AlarmStepAutoInstructionComplete, types.AlarmStepAutoInstructionFail:
		set["kpi_executed_auto_instructions"] = bson.M{"$setUnion": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$kpi_executed_auto_instructions",
				"then": "$kpi_executed_auto_instructions",
				"else": bson.A{},
			}},
			bson.A{event.Parameters.Instruction},
		}}
	}

	alarm := types.Alarm{}
	var err error
	if len(set) == 0 {
		err = p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
	} else {
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = p.alarmCollection.FindOneAndUpdate(ctx, match, []bson.M{{"$set": set}}, opts).Decode(&alarm)
	}

	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = alarmChangeType
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *instructionProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	p.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	sendTriggerEvent(ctx, event, result, p.amqpPublisher, p.encoder, p.logger)
}
