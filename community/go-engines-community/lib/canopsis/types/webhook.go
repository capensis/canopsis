package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

type WebhookRequest struct {
	URL        string            `bson:"url" json:"url" binding:"required,url"`
	Method     string            `bson:"method" json:"method" binding:"required"`
	Auth       *WebhookBasicAuth `bson:"auth,omitempty" json:"auth,omitempty"`
	Headers    map[string]string `bson:"headers,omitempty" json:"headers,omitempty"`
	Payload    string            `bson:"payload,omitempty" json:"payload,omitempty"`
	SkipVerify bool              `bson:"skip_verify" json:"skip_verify"`
}

type WebhookBasicAuth struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type WebhookDeclareTicket struct {
	EmptyResponse bool              `bson:"empty_response" json:"empty_response"`
	TicketID      string            `bson:"ticket_id" json:"ticket_id"`
	IsRegexp      bool              `bson:"is_regexp" json:"is_regexp"`
	Fields        map[string]string `bson:",inline"`
}

func (t *WebhookDeclareTicket) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
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

	fields := make(map[string]string)
	for k, v := range m {
		if strVal, ok := v.(string); ok {
			if k == "ticket_id" {
				t.TicketID = strVal
			} else {
				fields[k] = strVal
			}
		} else {
			return fmt.Errorf("invalid type of %s", k)
		}
	}
	t.Fields = fields

	return nil
}

func (t WebhookDeclareTicket) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"empty_response": t.EmptyResponse,
		"is_regexp":      t.IsRegexp,
		"ticket_id":      t.TicketID,
	}

	for k, v := range t.Fields {
		m[k] = v
	}

	return json.Marshal(m)
}
