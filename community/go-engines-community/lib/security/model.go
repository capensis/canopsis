package security

// User represents user model.
type User struct {
	ID             string
	Name           string
	Firstname      string
	Lastname       string
	Email          string
	HashedPassword string
	AuthApiKey     string
	Role           string
	Contact        struct {
		Name    string
		Address string
	}
	IsEnabled  bool
	ExternalID string
	Source     Source
}

type Source string

const SourceLdap Source = "ldap"
const SourceCas Source = "cas"
const SourceSaml Source = "saml"
