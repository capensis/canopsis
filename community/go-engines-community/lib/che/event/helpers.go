package event

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func updateMetaAlarmInfos(
	ctx context.Context,
	entityID string,
	newInfos map[string]types.Info,
	alarmCollection, entityCollection mongo.DbCollection,
) error {
	if len(newInfos) == 0 {
		return nil
	}

	cursor, err := alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"d":          entityID,
			"v.resolved": nil,
			"v.parents":  bson.M{"$nin": bson.A{nil, bson.A{}}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "v.parents",
			"foreignField": "d",
			"as":           "parents",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"v.resolved": nil,
					"cinfos":     bson.M{"$nin": bson.A{nil, bson.A{}}},
				}},
			},
		}},
		{"$unwind": "$parents"},
		{"$project": bson.M{
			"_id":    "$parents.d",
			"cinfos": "$parents.cinfos",
		}},
	})
	if err != nil {
		return fmt.Errorf("cannot find parents: %w", err)
	}

	defer cursor.Close(ctx)
	writeModels := make([]mongodriver.WriteModel, 0)
	for cursor.Next(ctx) {
		parent := struct {
			ID                      string   `bson:"_id"`
			EntityInfosFromChildren []string `bson:"cinfos"`
		}{}
		err = cursor.Decode(&parent)
		if err != nil {
			return fmt.Errorf("cannot decode parent: %w", err)
		}

		update := bson.M{}
		for _, infoName := range parent.EntityInfosFromChildren {
			if info, ok := newInfos[infoName]; ok {
				update["infos."+infoName+".value"] = info.Value
			}
		}

		if len(update) > 0 {
			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": parent.ID}).
				SetUpdate(bson.M{"$set": update}))
		}
	}

	if len(writeModels) > 0 {
		_, err = entityCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return fmt.Errorf("cannot update parent infos: %w", err)
		}
	}

	return nil
}

// runEventFilters runs the event filter service for the event,
// records its metrics on res and handles the asynchronous external-data suspension uniformly across processors.
// The returned bool is true when processing was suspended to fetch external data:
// res.EventFilterResult is then populated for the resume and the caller must stop and return res as is.
// The input result carries the accumulated state when an event is being resumed.
func runEventFilters(
	ctx context.Context,
	service eventfilter.Service,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	event *types.Event,
	res *ProcessorResult,
	partialRes *ProcessorResult,
) (eventfilter.ServiceResult, bool, error) {
	input := eventfilter.ServiceResult{}
	var preFilterEvent *types.Event
	if partialRes != nil {
		input = partialRes.EventFilterResult
		preFilterEvent = partialRes.PreFilterEvent
	} else {
		b, err := encoder.Encode(event)
		if err != nil {
			return eventfilter.ServiceResult{}, false, err
		}

		err = decoder.Decode(b, &preFilterEvent)
		if err != nil {
			return eventfilter.ServiceResult{}, false, err
		}
	}

	efr, err := service.ProcessEvent(ctx, event, input)
	if err != nil {
		return efr, false, err
	}

	res.EventMetric.ExecutedEnrichRules = efr.ExecutedEnrichRuleCount
	res.EventMetric.ExternalRequests = efr.ExternalRequestCount

	if efr.ExternalDataRequest != nil {
		res.EventFilterResult = efr
		res.PreFilterEvent = preFilterEvent

		return efr, true, nil
	}

	return efr, false, nil
}

func logInfosUpdate(metricsSender metrics.EntityInfosUpdateSender, entityID string, updatedInfos map[string]eventfilter.UpdatedValue) {
	now := time.Now()
	for k, v := range updatedInfos {
		metricsSender.Send(now, entityID, v.RuleID, k, v.NewValue)
	}
}
