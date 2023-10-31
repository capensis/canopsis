// Package mongoadapter contains casbin mongo adapter.
// Adapter loads policy from mongo collection and transform result into casbin format.
// Refactor mongo collection and use casbin mongo adapter after all API is migrated to Go.
package mongoadapter

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libmodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	casbinPtypePolicy = "p"
	casbinPtypeRole   = "g"
)

// NewAdapter creates mongo adapter.
func NewAdapter(db mongo.DbClient) persist.Adapter {
	return &adapter{
		userCollection:       db.Collection(mongo.UserCollection),
		roleCollection:       db.Collection(mongo.RoleCollection),
		permissionCollection: db.Collection(mongo.PermissionCollection),
	}
}

// adapter implements casbin adapter interface.
type adapter struct {
	userCollection       mongo.DbCollection
	roleCollection       mongo.DbCollection
	permissionCollection mongo.DbCollection
}

type role struct {
	ID          string                `bson:"_id"`
	Permissions map[string]permission `bson:"permissions"`
}

type permission struct {
	Bitmask int64  `bson:"bitmask"`
	Type    string `bson:"type"`
}

type user struct {
	ID    string   `bson:"_id"`
	Roles []string `bson:"roles"`
}

// LoadPolicy loads all policy rules from mongo collection.
func (a *adapter) LoadPolicy(model model.Model) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := a.loadRoles(ctx, model)
	if err != nil {
		return err
	}

	err = a.loadSubjects(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

// SavePolicy isn't implemented.
// Implement it if all API is migrated to Go
// and there is not time to refactor mongo collection.
func (*adapter) SavePolicy(model.Model) error {
	panic("implement me")
}

func (*adapter) AddPolicy(string, string, []string) error {
	panic("implement me")
}

func (*adapter) RemovePolicy(string, string, []string) error {
	panic("implement me")
}

func (*adapter) RemoveFilteredPolicy(string, string, int, ...string) error {
	panic("implement me")
}

// loadRoles fetches roles from mongo collection and adds them to casbin policy.
func (a *adapter) loadRoles(
	ctx context.Context,
	model model.Model,
) (resErr error) {
	roles, err := a.getRolesWithPermissions(ctx)
	if err != nil {
		return err
	}

	ptype := casbinPtypePolicy
	sec := ptype
	permBitmasksByType := map[string]map[string]int64{
		libmodel.ObjectTypeCRUD: {
			libmodel.PermissionCreate: libmodel.PermissionBitmaskCreate,
			libmodel.PermissionRead:   libmodel.PermissionBitmaskRead,
			libmodel.PermissionUpdate: libmodel.PermissionBitmaskUpdate,
			libmodel.PermissionDelete: libmodel.PermissionBitmaskDelete,
		},
		libmodel.ObjectTypeRW: {
			libmodel.PermissionRead:   libmodel.PermissionBitmaskRead,
			libmodel.PermissionUpdate: libmodel.PermissionBitmaskUpdate,
			libmodel.PermissionDelete: libmodel.PermissionBitmaskDelete,
		},
	}

	for _, r := range roles {
		for objId, obj := range r.Permissions {
			if obj.Type != "" {
				if permBitmasksByName, ok := permBitmasksByType[obj.Type]; ok {
					for permName, bitmask := range permBitmasksByName {
						if obj.Bitmask&bitmask == bitmask {
							model.AddPolicy(sec, ptype, []string{r.ID, objId, permName})
						}
					}
				} else {
					return fmt.Errorf("unknown config type \"%s\"", obj.Type)
				}
			} else if obj.Bitmask&libmodel.PermissionBitmaskCan == libmodel.PermissionBitmaskCan {
				model.AddPolicy(sec, ptype, []string{r.ID, objId, libmodel.PermissionCan})
			}
		}
	}

	return nil
}

// loadSubjects loads subjects from mongo collection and adds them to casbin policy.
func (a *adapter) loadSubjects(ctx context.Context, model model.Model) (resErr error) {
	cursor, err := a.userCollection.Find(ctx, bson.M{"roles": bson.M{"$exists": true, "$ne": ""}})
	if err != nil {
		return err
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil && resErr == nil {
			resErr = err
		}
	}()

	ptype := casbinPtypeRole
	sec := ptype

	for cursor.Next(ctx) {
		var u user
		err := cursor.Decode(&u)
		if err != nil {
			return err
		}

		for _, r := range u.Roles {
			model.AddPolicy(sec, ptype, []string{u.ID, r})
		}
	}

	return nil
}

type permissionDocument struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Type string `bson:"type"`
}

type roleDocument struct {
	ID          string           `bson:"_id"`
	Permissions map[string]int64 `bson:"permissions"`
}

func (a *adapter) getRolesWithPermissions(ctx context.Context) ([]role, error) {
	// Read all permissions into memory
	var permissions []permissionDocument
	cursor, err := a.permissionCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &permissions)
	if err != nil {
		return nil, err
	}

	// Create a map of permission IDs to permission names
	permissionByIDs := make(map[string]permissionDocument)
	for _, permission := range permissions {
		permissionByIDs[permission.ID] = permission
	}

	// Read all roles into memory
	var roles []roleDocument
	cursorRole, err := a.roleCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursorRole.All(ctx, &roles)
	if err != nil {
		return nil, err
	}

	// Create a slice of role type
	var result []role
	for _, roleItem := range roles {
		permissions := make(map[string]permission)
		for permissionID, mask := range roleItem.Permissions {
			permissionItem := permissionByIDs[permissionID]
			permissions[permissionItem.Name] = permission{
				Bitmask: mask,
				Type:    permissionItem.Type,
			}
		}
		result = append(result, role{
			ID:          roleItem.ID,
			Permissions: permissions,
		})
	}

	return result, nil
}
