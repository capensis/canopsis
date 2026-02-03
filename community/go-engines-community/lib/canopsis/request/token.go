package request

type WebhookAuthToken struct {
	Rule         string `bson:"rule" json:"rule"`
	QueryString  string `bson:"query_string,omitempty" json:"query_string,omitempty"`
	Header       string `bson:"header,omitempty" json:"header,omitempty"`
	PayloadField string `bson:"payload_field,omitempty" json:"payload_field,omitempty"`
}
