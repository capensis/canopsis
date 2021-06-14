package ics

import "time"

// Event is used to implement ICS event component.
type Event interface {
	serializable
	SetUID(string)
	SetCreated(time.Time)
	SetLastModified(time.Time)
	SetDTStart(time.Time)
	SetDTEnd(time.Time)
	SetExdate(t time.Time)
	SetDTStamp(time.Time)
	SetRRule(string)
	SetDescription(string)
	SetSummary(string)
	SetTransp(string)
	SetContact(string)
	SetPriority(int64)
}

// NewEvent creates new event.
func NewEvent() Event {
	return &event{
		component: component{
			ComponentType: "VEVENT",
		},
	}
}

// event represents ICS event component.
type event struct {
	component
}

func (e *event) SetUID(s string) {
	e.AddProperty("UID", s)
}

func (e *event) SetCreated(t time.Time) {
	e.AddDateTimeProperty("CREATED", t)
}

func (e *event) SetLastModified(t time.Time) {
	e.AddDateTimeProperty("LAST-MODIFIED", t)
}

func (e *event) SetDTStart(t time.Time) {
	e.AddDateTimeProperty("DTSTART", t)
}

func (e *event) SetDTEnd(t time.Time) {
	e.AddDateTimeProperty("DTEND", t)
}

func (e *event) SetExdate(t time.Time) {
	e.AddDateTimeProperty("EXDATE", t)
}

func (e *event) SetDTStamp(t time.Time) {
	e.AddDateTimeProperty("DTSTAMP", t)
}

func (e *event) SetRRule(s string) {
	e.AddProperty("RRULE", s)
}

func (e *event) SetDescription(s string) {
	e.AddProperty("DESCRIPTION", s)
}

func (e *event) SetSummary(s string) {
	e.AddProperty("SUMMARY", s)
}

func (e *event) SetTransp(s string) {
	e.AddProperty("TRANSP", s)
}

func (e *event) SetContact(s string) {
	e.AddProperty("CONTACT", s)
}

func (e *event) SetPriority(s int64) {
	e.AddProperty("PRIORITY", s)
}
