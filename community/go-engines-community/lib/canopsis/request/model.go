package request

//todo: copy from webhook package, webhook package should use this package instead of its own models

type Parameters struct {
	URL        string            `bson:"url,omitempty" json:"url,omitempty"`
	Method     string            `bson:"method,omitempty" json:"method,omitempty"`
	Auth       *BasicAuth        `bson:"auth,omitempty" json:"auth,omitempty"`
	Headers    map[string]string `bson:"headers,omitempty" json:"headers,omitempty"`
	Payload    string            `bson:"payload,omitempty" json:"payload,omitempty"`
	SkipVerify bool              `bson:"skip_verify" json:"skip_verify"`
}

type BasicAuth struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}
