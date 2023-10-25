package security

const (
	SourceLdap = "ldap"
	SourceCas  = "cas"
	SourceSaml = "saml"
)

// User represents user model.
type User struct {
	ID             string   `bson:"_id"`
	Name           string   `bson:"name"`
	DisplayName    string   `bson:"display_name,omitempty"`
	Firstname      string   `bson:"firstname"`
	Lastname       string   `bson:"lastname"`
	Email          string   `bson:"email"`
	HashedPassword string   `bson:"password,omitempty"`
	AuthApiKey     string   `bson:"authkey"`
	Roles          []string `bson:"roles"`
	Contact        struct {
		Name    string `bson:"name"`
		Address string `bson:"address"`
	} `bson:"contact"`
	IsEnabled  bool   `bson:"enable"`
	ExternalID string `bson:"external_id"`
	Source     string `bson:"source"`
}
