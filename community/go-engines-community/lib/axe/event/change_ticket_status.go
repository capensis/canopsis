package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewChangeTicketStatusProcessor(client mongo.DbClient) Processor {
	return &changeTicketStatusProcessor{
		dbClient:                client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
	}
}

type changeTicketStatusProcessor struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
}

func (p *changeTicketStatusProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	if event.Parameters.Ticket == "" {
		return Result{}, nil
	}

	match := bson.M{
		"_id": bson.M{"$in": event.AlarmIDs},
		"v.tickets": bson.M{
			"$elemMatch": bson.M{
				"ticket":             event.Parameters.Ticket,
				"ticket_system_name": event.Parameters.TicketSystemName,
				// to prevent status change after closed, closed should be the last status
				"ticket_status": bson.M{"$ne": webhook.TicketStatusClosed},
				"$or": bson.A{
					bson.M{"ticket_status": bson.M{"$ne": event.Parameters.TicketStatus}},
					bson.M{"ticket_source_status": bson.M{"$ne": event.Parameters.TicketSourceStatus}},
				},
			},
		},
	}

	update := []bson.M{
		{"$set": bson.M{
			"v.steps": bson.M{
				"$cond": bson.A{
					bson.M{"$and": bson.A{
						bson.M{"$lt": bson.A{
							bson.M{"$size": "$v.steps"},
							types.AlarmStepsHardLimit,
						}},
						bson.M{"$gt": bson.A{
							bson.M{"$size": bson.M{
								"$filter": bson.M{
									"input": "$v.tickets",
									"as":    "ticket",
									"cond": bson.M{"$and": bson.A{
										bson.M{"$eq": bson.A{"$$ticket.ticket", event.Parameters.Ticket}},
										bson.M{"$eq": bson.A{"$$ticket.ticket_system_name", event.Parameters.TicketSystemName}},
										bson.M{"$ne": bson.A{"$$ticket.ticket_status", event.Parameters.TicketStatus}},
									}},
								},
							}},
							0,
						}},
					}},
					addStepUpdateQuery(
						valStepUpdateQueryWithInPbhInterval(
							types.AlarmStepChangeTicketStatus,
							types.CpsNumber(event.Parameters.TicketStatus),
							event.Parameters.Output,
							event.Parameters,
						),
					),
					"$v.steps",
				},
			},
			"v.last_update_date": event.Parameters.Timestamp,
			"v.tickets": bson.M{
				"$map": bson.M{
					"input": "$v.tickets",
					"as":    "ticket",
					"in": bson.M{
						"$cond": bson.A{
							bson.M{"$and": bson.A{
								bson.M{"$eq": bson.A{"$$ticket.ticket", event.Parameters.Ticket}},
								bson.M{"$eq": bson.A{"$$ticket.ticket_system_name", event.Parameters.TicketSystemName}},
								bson.M{"$or": bson.A{
									bson.M{"$ne": bson.A{"$$ticket.ticket_status", event.Parameters.TicketStatus}},
									bson.M{"$ne": bson.A{"$$ticket.ticket_source_status", event.Parameters.TicketSourceStatus}},
								}},
							}},
							bson.M{"$mergeObjects": bson.A{"$$ticket", bson.M{
								"ticket_status":             event.Parameters.TicketStatus,
								"ticket_prev_status":        event.Parameters.TicketPrevStatus,
								"ticket_source_status":      event.Parameters.TicketSourceStatus,
								"ticket_prev_source_status": event.Parameters.TicketPrevSourceStatus,
								"ticket_last_check_time":    event.Parameters.TicketLastCheckTime,
							}}},
							"$$ticket",
						},
					},
				},
			},
			"v.ticket": bson.M{
				"$cond": bson.A{
					bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$v.ticket.ticket", event.Parameters.Ticket}},
						bson.M{"$eq": bson.A{"$v.ticket.ticket_system_name", event.Parameters.TicketSystemName}},
						bson.M{"$or": bson.A{
							bson.M{"$ne": bson.A{"$v.ticket.ticket_status", event.Parameters.TicketStatus}},
							bson.M{"$ne": bson.A{"$v.ticket.ticket_source_status", event.Parameters.TicketSourceStatus}},
						}},
					}},
					bson.M{"$mergeObjects": bson.A{"$v.ticket", bson.M{
						"ticket_status":             event.Parameters.TicketStatus,
						"ticket_prev_status":        event.Parameters.TicketPrevStatus,
						"ticket_source_status":      event.Parameters.TicketSourceStatus,
						"ticket_prev_source_status": event.Parameters.TicketPrevSourceStatus,
						"ticket_last_check_time":    event.Parameters.TicketLastCheckTime,
					}}},
					"$v.ticket",
				},
			},
		}},
	}

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := p.alarmCollection.UpdateMany(ctx, match, update)
		if err != nil {
			return fmt.Errorf("failed to update ticket status to opened alarms: %w", err)
		}

		_, err = p.resolvedAlarmCollection.UpdateMany(ctx, match, update)
		if err != nil {
			return fmt.Errorf("failed to update ticket status to resolved alarms: %w", err)
		}

		return nil
	})
	if err != nil {
		return Result{}, err
	}

	return Result{}, nil
}
