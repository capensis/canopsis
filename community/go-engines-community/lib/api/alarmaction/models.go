package alarmaction

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

type Request struct {
	Comment string `json:"comment" binding:"max=255"`
}

type AckRequest struct {
	AckResources bool   `json:"ack_resources"`
	Comment      string `json:"comment" binding:"max=255"`
}

type SnoozeRequest struct {
	Duration types.DurationWithUnit `json:"duration" binding:"required"`
	Comment  string                 `json:"comment" binding:"max=255"`
}

type AssocTicketRequest struct {
	Ticket          string            `json:"ticket" binding:"required"`
	Url             string            `json:"url" binding:"max=255"`
	Comment         string            `json:"comment" binding:"max=255"`
	SystemName      string            `json:"system_name" binding:"max=255"`
	Data            map[string]string `json:"data"`
	TicketResources bool              `json:"ticket_resources"`
}

type ChangeStateRequest struct {
	State   *int64 `json:"state" binding:"required,oneof=0 1 2 3"`
	Comment string `json:"comment" binding:"max=255"`
}

type CommentRequest struct {
	Comment string `json:"comment" binding:"required,max=255"`
}
