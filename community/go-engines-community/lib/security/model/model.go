package model

// Rbac represents mongo collection structure.
// Collection contains
// - objects with defined ID and Name fields
// - roles with defined ID, Name and PermConfigList fields
// - subjects with defined ID, Role and data fields
type Rbac struct {
	ID   string   `bson:"_id,omitempty"`
	Type LineType `bson:"crecord_type"`
	Name string   `bson:"crecord_name"`
	// ObjectType defines if object has can permission or CRUD permissions.
	ObjectType string `bson:"type,omitempty"`
	// Role is only for subject.
	Role string `bson:"role,omitempty"`
	// PermConfigList is only for role.
	PermConfigList map[string]struct {
		// Bitmask contains bitmask of CRUD permissions.
		Bitmask int `bson:"checksum"`
	} `bson:"rights,omitempty"`
	// Following fields contains extra data for subject.
	Email          string `bson:"mail,omitempty"`
	Lastname       string `bson:"lastname,omitempty"`
	Firstname      string `bson:"firstname,omitempty"`
	HashedPassword string `bson:"shadowpasswd,omitempty"`
	AuthApiKey     string `bson:"authkey,omitempty"`
	IsEnabled      bool   `bson:"enable,omitempty"`
	ExternalID     string `bson:"external_id,omitempty"`
	Source         string `bson:"source,omitempty"`
	Contact        struct {
		Name    string `bson:"name,omitempty"`
		Address string `bson:"address,omitempty"`
	} `bson:"contact,omitempty"`
}

type LineType string

const (
	LineTypeObject  LineType = "action"
	LineTypeRole    LineType = "role"
	LineTypeSubject LineType = "user"
)

const LineObjectTypeCRUD = "CRUD"
