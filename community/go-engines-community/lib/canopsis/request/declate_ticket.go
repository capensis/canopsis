package request

import (
	"encoding/json"
	"errors"
	"fmt"
)

type WebhookDeclareTicket struct {
	EmptyResponse bool              `bson:"empty_response" json:"empty_response"`
	IsRegexp      bool              `bson:"is_regexp" json:"is_regexp"`
	TicketID      string            `bson:"ticket_id,omitempty" json:"ticket_id"`
	TicketIdTpl   string            `bson:"ticket_id_tpl,omitempty" json:"ticket_id_tpl"`
	TicketUrl     string            `bson:"ticket_url,omitempty" json:"ticket_url"`
	TicketUrlTpl  string            `bson:"ticket_url_tpl,omitempty" json:"ticket_url_tpl"`
	CustomFields  map[string]string `bson:",inline"`
}

func (t *WebhookDeclareTicket) UnmarshalJSON(b []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if emptyResponse, ok := m["empty_response"]; ok {
		if boolVal, ok := emptyResponse.(bool); ok {
			t.EmptyResponse = boolVal
			delete(m, "empty_response")
		} else {
			return errors.New("invalid type of empty_response")
		}
	}

	if isRegexp, ok := m["is_regexp"]; ok {
		if boolVal, ok := isRegexp.(bool); ok {
			t.IsRegexp = boolVal
			delete(m, "is_regexp")
		} else {
			return errors.New("invalid type of is_regexp")
		}
	}

	customFields := make(map[string]string)
	for k, v := range m {
		if strVal, ok := v.(string); ok {
			switch k {
			case "ticket_id":
				t.TicketID = strVal
			case "ticket_id_tpl":
				t.TicketIdTpl = strVal
			case "ticket_url":
				t.TicketUrl = strVal
			case "ticket_url_tpl":
				t.TicketUrlTpl = strVal
			default:
				customFields[k] = strVal
			}
		} else {
			return fmt.Errorf("invalid type of %s", k)
		}
	}
	t.CustomFields = customFields

	return nil
}

func (t WebhookDeclareTicket) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"empty_response": t.EmptyResponse,
		"is_regexp":      t.IsRegexp,
		"ticket_id":      t.TicketID,
		"ticket_id_tpl":  t.TicketIdTpl,
		"ticket_url":     t.TicketUrl,
		"ticket_url_tpl": t.TicketUrlTpl,
	}

	for k, v := range t.CustomFields {
		m[k] = v
	}

	return json.Marshal(m)
}
