package security_test

import (
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/kylelemons/godebug/pretty"
)

func TestUser_SetRolesFromIdp(t *testing.T) {
	dataSets := []struct {
		user          security.User
		newIdpRoles   []string
		expectedRoles []string
		mergeRoles    bool
	}{
		{
			user:          security.User{},
			newIdpRoles:   []string{"role-1", "role-2"},
			expectedRoles: []string{"role-1", "role-2"},
			mergeRoles:    false,
		},
		{
			user: security.User{
				Roles: []string{"role-1", "role-2"},
			},
			newIdpRoles:   []string{"role-2", "role-3"},
			expectedRoles: []string{"role-2", "role-3"},
			mergeRoles:    false,
		},
		{
			user: security.User{
				Roles: []string{"role-1", "role-2"},
			},
			newIdpRoles:   []string{},
			expectedRoles: []string{},
			mergeRoles:    false,
		},
		{
			user: security.User{
				Roles:    []string{"role-1", "role-2"},
				IdPRoles: []string{"role-1", "role-2"},
			},
			newIdpRoles:   []string{},
			expectedRoles: []string{},
			mergeRoles:    false,
		},
		{
			user:          security.User{},
			newIdpRoles:   []string{"role-1", "role-2"},
			expectedRoles: []string{"role-1", "role-2"},
			mergeRoles:    true,
		},
		{
			user: security.User{
				Roles: []string{"role-1", "role-2"},
			},
			newIdpRoles:   []string{"role-2", "role-3"},
			expectedRoles: []string{"role-2", "role-3", "role-1"},
			mergeRoles:    true,
		},
		{
			user: security.User{
				Roles:    []string{"role-1", "role-2", "role-3"},
				IdPRoles: []string{"role-1", "role-2"},
			},
			newIdpRoles:   []string{"role-1"},
			expectedRoles: []string{"role-1", "role-3"},
			mergeRoles:    true,
		},
		{
			user: security.User{
				Roles:    []string{"role-1", "role-2", "role-3"},
				IdPRoles: []string{"role-1", "role-3"},
			},
			newIdpRoles:   []string{},
			expectedRoles: []string{"role-2"},
			mergeRoles:    true,
		},
	}

	for i, set := range dataSets {
		set.user.SetRolesFromIdP(set.newIdpRoles, set.mergeRoles)
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if diff := pretty.Compare(set.expectedRoles, set.user.Roles); diff != "" {
				t.Errorf("unexpected roles: %s", diff)
			}

			if diff := pretty.Compare(set.newIdpRoles, set.user.IdPRoles); diff != "" {
				t.Errorf("unexpected idp roles: %s", diff)
			}
		})
	}
}
