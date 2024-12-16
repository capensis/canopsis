package security

const (
	SourceLdap   = "ldap"
	SourceCas    = "cas"
	SourceSaml   = "saml"
	SourceOauth2 = "oauth2"
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

	// IdpRoles field show roles from idp, and they should be used ONLY in idp/canopsis role merging, see SetRolesFromIdp.
	IdpRoles []string `bson:"idp_roles"`
}

func (u *User) SetRolesFromIdp(newIdpRoles []string, mergeRoles bool) {
	if mergeRoles {
		oldIdpRolesMap := make(map[string]bool)
		for _, role := range u.IdpRoles {
			oldIdpRolesMap[role] = true
		}

		newIdpRolesMap := make(map[string]bool)
		for _, role := range newIdpRoles {
			newIdpRolesMap[role] = true
		}

		u.IdpRoles = make([]string, len(newIdpRoles))
		copy(u.IdpRoles, newIdpRoles)

		for _, role := range u.Roles {
			if !oldIdpRolesMap[role] && !newIdpRolesMap[role] {
				newIdpRoles = append(newIdpRoles, role)
			}
		}

		u.Roles = newIdpRoles
	} else {
		u.IdpRoles = newIdpRoles
		u.Roles = newIdpRoles
	}
}
