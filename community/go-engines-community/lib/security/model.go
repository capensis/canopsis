package security

// User represents user model.
type User struct {
	ID             string   `bson:"_id"`
	Name           string   `bson:"name"`
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
	Source     Source `bson:"source"`
}

type Source string

const SourceLdap Source = "ldap"
const SourceCas Source = "cas"
const SourceSaml Source = "saml"
