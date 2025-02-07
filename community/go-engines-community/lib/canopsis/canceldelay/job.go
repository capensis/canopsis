package canceldelay

type Job struct {
	ID        string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string `bson:"name" json:"name"`
	Component string `bson:"comp" json:"comp"`
	Type      string `bson:"type" json:"type"`
	Delay     int64  `bson:"delay" json:"delay"`
	ExecTime  int64  `bson:"exec_time" json:"exec_time"`
	ResendAt  int64  `bson:"resend_at,omitempty" json:"resend_at,omitempty"`
}
