// ics contains ics calendar implementation.
package ics

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

const timeFormat = "20060102T150405Z"
const (
	UndefinedPriority = 0
	MinPriority       = 1
	MaxPriority       = 9
)

// Calendar is used to implement ICS Calendar structure.
type Calendar interface {
	// AddComponent adds ICS component to Calendar.
	AddComponent(component serializable)
	// Read returns Calendar content.
	Read() (io.Reader, int64)
}

// NewCalendar creates new calendar.
func NewCalendar() Calendar {
	return &calendar{
		Version: "2.0",
	}
}

// serializable is used to implement ICS elements saving.
type serializable interface {
	serialize(writer io.Writer)
}

// calendar represents ICS calendar.
type calendar struct {
	Version    string
	Components []serializable
}

// AddComponent adds component to slice.
func (c *calendar) AddComponent(component serializable) {
	if c.Components == nil {
		c.Components = make([]serializable, 0)
	}

	c.Components = append(c.Components, component)
}

// Read serializes calendar content to string buffer.
func (c *calendar) Read() (io.Reader, int64) {
	b := bytes.NewBufferString("")
	c.serialize(b)

	return b, int64(b.Len())
}

// serialize saves calendar components and calendar properties.
func (c *calendar) serialize(w io.Writer) {
	_, _ = fmt.Fprint(w, "BEGIN:VCALENDAR\n")
	_, _ = fmt.Fprintf(w, "VERSION:%s\n", c.Version)

	for _, c := range c.Components {
		c.serialize(w)
	}

	_, _ = fmt.Fprint(w, "END:VCALENDAR\n")
}

// component represents ICS calendar component.
type component struct {
	ComponentType string
	Properties    []serializable
}

// addProperty adds property to slice.
func (c *component) addProperty(p serializable) {
	if c.Properties == nil {
		c.Properties = make([]serializable, 0)
	}

	c.Properties = append(c.Properties, p)
}

// AddProperty adds simple property.
func (c *component) AddProperty(name string, val interface{}) {
	c.addProperty(&property{
		propertyName: name,
		value:        val,
	})
}

// AddDateTimeProperty adds date-time property.
func (c *component) AddDateTimeProperty(name string, val time.Time) {
	c.addProperty(&dateTimeProperty{
		propertyName: name,
		value:        val,
	})
}

// serialize saves component properties.
func (c *component) serialize(w io.Writer) {
	_, _ = fmt.Fprintf(w, "BEGIN:%s\n", c.ComponentType)
	for _, p := range c.Properties {
		p.serialize(w)
	}
	_, _ = fmt.Fprintf(w, "END:%s\n", c.ComponentType)
}

// property represents ICS calendar simple property.
type property struct {
	propertyName string
	value        interface{}
}

// serialize saves property value.
func (p *property) serialize(w io.Writer) {
	if p.value != "" {
		_, _ = fmt.Fprintf(w, "%s:%v\n", p.propertyName, p.value)
	}
}

// property represents ICS calendar date-time property.
type dateTimeProperty struct {
	propertyName string
	value        time.Time
}

// serialize saves property value in required date-time format.
func (p *dateTimeProperty) serialize(w io.Writer) {
	if !p.value.IsZero() {
		if p.value.Location() == time.UTC {
			_, _ = fmt.Fprintf(w, "%s:%s\n", p.propertyName, p.value.Format(timeFormat))
		} else {
			_, _ = fmt.Fprintf(
				w,
				"%s;TZID=%s:%s\n",
				p.propertyName,
				p.value.Location().String(),
				p.value.Format(timeFormat),
			)
		}
	}
}
