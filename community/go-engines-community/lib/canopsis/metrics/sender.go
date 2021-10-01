package metrics

import (
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

const DefaultSlice = "default_slice"

const TotalAlarmNumberEntity = "total_alarm_number_entity"
const TotalAlarmNumberSlice = "total_alarm_number_slice"

const PbhAlarmNumberEntity = "pbh_alarm_number_entity"
const PbhAlarmNumberSlice = "pbh_alarm_number_slice"

const InstructionAlarmNumberEntity = "instruction_alarm_number_entity"
const InstructionAlarmNumberSlice = "instruction_alarm_number_slice"

const TicketAlarmNumberEntity = "ticket_alarm_number_entity"
const TicketAlarmNumberSlice = "ticket_alarm_number_slice"

const CorrelationAlarmNumberEntity = "correlation_alarm_number_entity"
const CorrelationAlarmNumberSlice = "correlation_alarm_number_slice"

const AckAlarmNumberEntity = "ack_alarm_number_entity"
const AckAlarmNumberSlice = "ack_alarm_number_slice"

const CancelAckAlarmNumberEntity = "cancel_ack_alarm_number_entity"
const CancelAckAlarmNumberSlice = "cancel_ack_alarm_number_slice"

const AckDurationEntity = "ack_duration_entity"
const AckDurationSlice = "ack_duration_slice"

const ResolveDurationEntity = "resolve_duration_entity"
const ResolveDurationSlice = "resolve_duration_slice"

type Parameters struct {
	Value float64 `json:"value"`
}

type Event struct {
	Type       string            `json:"type"`
	Labels     prometheus.Labels `json:"labels"`
	Parameters Parameters        `json:"parameters"`
}

type Sender interface {
	HandleMetricsByEvent(event types.Event)
	HandleMetricForMetaalarmChild(child types.AlarmWithEntity)
}

type sender struct {
	pubExchangeName string
	pubQueueName    string
	pubChannel      libamqp.Channel
	encoder         encoding.Encoder
	logger          zerolog.Logger
}

func NewSender(
	pubExchangeName,
	pubQueueName string,
	pubChannel libamqp.Channel,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Sender {
	return &sender{
		pubExchangeName: pubExchangeName,
		pubQueueName:    pubQueueName,
		pubChannel:      pubChannel,
		encoder:         encoder,
		logger:          logger,
	}
}

//TODO: for every slice metrics: check if belongs to a slice! Only default for now.
func (s *sender) HandleMetricsByEvent(event types.Event) {
	alarmChange := event.AlarmChange
	alarm := event.Alarm
	entity := event.Entity

	if alarm == nil || entity == nil || alarmChange == nil {
		return
	}
	
	entityID := alarm.EntityID
	category := entity.Category
	author := event.Author

	switch alarmChange.Type {
	case types.AlarmChangeTypeCreate:
		s.sendMetric(Event{
			Type:   TotalAlarmNumberEntity,
			Labels: prometheus.Labels{"entityID": entityID, "category": category},
		})

		s.sendMetric(Event{
			Type:   TotalAlarmNumberSlice,
			Labels: prometheus.Labels{"slice": DefaultSlice},
		})
	case types.AlarmChangeTypeCreateAndPbhEnter:
		s.sendMetric(Event{
			Type:   PbhAlarmNumberEntity,
			Labels: prometheus.Labels{"entityID": entityID, "category": category},
		})

		s.sendMetric(Event{
			Type:   PbhAlarmNumberSlice,
			Labels: prometheus.Labels{"slice": DefaultSlice},
		})
	case types.AlarmChangeTypeAssocTicket:
		foundStep := false
		for _, step := range alarm.Value.Steps {
			if step.Type == types.AlarmStepAssocTicket {
				// if it's already had an assoc ticket, then don't send metric
				if foundStep {
					return
				}

				foundStep = true
			}
		}

		s.sendMetric(Event{
			Type:   TicketAlarmNumberEntity,
			Labels: prometheus.Labels{"entityID": entityID, "category": category, "username": author},
		})

		s.sendMetric(Event{
			Type:   TicketAlarmNumberSlice,
			Labels: prometheus.Labels{"slice": DefaultSlice, "username": author},
		})
	case types.AlarmChangeTypeAck:
		if alarm.Value.ACK == nil {
			return
		}

		s.sendMetric(Event{
			Type:   AckAlarmNumberEntity,
			Labels: prometheus.Labels{"entityID": entityID, "category": category, "username": author},
		})

		s.sendMetric(Event{
			Type:   AckAlarmNumberSlice,
			Labels: prometheus.Labels{"slice": DefaultSlice, "username": author},
		})

		ackDuration := alarm.Value.ACK.Timestamp.Sub(alarm.Value.CreationDate.Time).Seconds()
		s.sendMetric(Event{
			Type:       AckDurationEntity,
			Labels:     prometheus.Labels{"entityID": entityID, "category": category, "username": author},
			Parameters: Parameters{Value: ackDuration},
		})

		s.sendMetric(Event{
			Type:       AckDurationSlice,
			Labels:     prometheus.Labels{"slice": DefaultSlice, "username": author},
			Parameters: Parameters{Value: ackDuration},
		})
	case types.AlarmChangeTypeAckremove:
		s.sendMetric(Event{
			Type:   CancelAckAlarmNumberEntity,
			Labels: prometheus.Labels{"entityID": entityID, "category": category, "username": author},
		})

		s.sendMetric(Event{
			Type:   CancelAckAlarmNumberSlice,
			Labels: prometheus.Labels{"slice": DefaultSlice, "username": author},
		})
	case types.AlarmChangeTypeAutoInstructionStart:
		foundStep := false
		for _, step := range alarm.Value.Steps {
			if step.Type == types.AlarmStepAutoInstructionStart {
				// if it's already had completed auto-instruction, then don't send metric
				if foundStep {
					return
				}

				foundStep = true
			}
		}

		s.sendMetric(Event{
			Type:   InstructionAlarmNumberEntity,
			Labels: prometheus.Labels{"entityID": entityID, "category": category},
		})

		s.sendMetric(Event{
			Type:   InstructionAlarmNumberSlice,
			Labels: prometheus.Labels{"slice": DefaultSlice},
		})
	case types.AlarmChangeTypeResolve:
		if alarm.Value.Resolved == nil {
			return
		}

		resolveDuration := alarm.Value.Resolved.Sub(alarm.Value.CreationDate.Time).Seconds()
		s.sendMetric(Event{
			Type:       ResolveDurationEntity,
			Labels:     prometheus.Labels{"entityID": entityID, "category": category},
			Parameters: Parameters{Value: resolveDuration},
		})

		s.sendMetric(Event{
			Type:       ResolveDurationSlice,
			Labels:     prometheus.Labels{"slice": DefaultSlice},
			Parameters: Parameters{Value: resolveDuration},
		})
	}
}

func (s *sender) HandleMetricForMetaalarmChild(child types.AlarmWithEntity) {
	s.sendMetric(Event{
		Type:   CorrelationAlarmNumberEntity,
		Labels: prometheus.Labels{"entityID": child.Entity.ID, "category": child.Entity.Category},
	})

	s.sendMetric(Event{
		Type:   CorrelationAlarmNumberSlice,
		Labels: prometheus.Labels{"slice": DefaultSlice},
	})
}

func (s *sender) sendMetric(event Event) {
	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to serialize metrics event", event.Type)
	}

	err = s.pubChannel.Publish(
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		amqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to send metrics event", event.Type)
	}
}
