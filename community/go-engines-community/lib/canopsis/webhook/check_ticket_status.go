package webhook

import (
	"context"
)

const (
	TicketStatusUnknown = iota
	TicketStatusOpen
	TicketStatusAssigned
	TicketStatusInProgress
	TicketStatusClosed
)

func GetStatusName(status int) string {
	switch status {
	case TicketStatusUnknown:
		return "Unknown"
	case TicketStatusOpen:
		return "Open"
	case TicketStatusInProgress:
		return "In Progress"
	case TicketStatusAssigned:
		return "Assigned"
	case TicketStatusClosed:
		return "Closed"
	default:
		return "Unknown"
	}
}

type CheckTicketStatusWorker interface {
	Run(ctx context.Context) error
}

type CheckTicketStatusService interface {
	CreateCheckStatusJob(ctx context.Context, historyID string, ticketStatus int, ticketSourceStatus string) (CheckTicketStatusJob, error)
	AddAlarmToCheckStatusJob(ctx context.Context, jobID, alarmID string) (CheckTicketStatusJob, error)
	RemoveAlarmFromCheckStatusJob(ctx context.Context, ticketID, ticketSystemName, alarmID string) error
}

type nullJobStatusService struct{}

func NewNullJobStatusService() CheckTicketStatusService {
	return &nullJobStatusService{}
}

func (nullJobStatusService) CreateCheckStatusJob(_ context.Context, _ string, _ int, _ string) (CheckTicketStatusJob, error) {
	return CheckTicketStatusJob{}, nil
}

func (nullJobStatusService) AddAlarmToCheckStatusJob(_ context.Context, _, _ string) (CheckTicketStatusJob, error) {
	return CheckTicketStatusJob{}, nil
}

func (nullJobStatusService) RemoveAlarmFromCheckStatusJob(_ context.Context, _, _, _ string) error {
	return nil
}
