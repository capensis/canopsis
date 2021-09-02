// mongoadapter contains casbin mongo adapter.
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

// NewAdapter creates mongo adapter.
func NewAdapter(db mongo.DbClient) persist.Adapter {
	return &adapter{
		collection: db.Collection(mongo.RightsMongoCollection),
	}
}

// adapter implements casbin adapter interface.
type adapter struct {
	collection mongo.DbCollection
}

const (
	CasbinPtypePolicy = "p"
	CasbinPtypeRole   = "g"
)

type objectConfig struct {
	Name string
	Type string
}

// LoadPolicy loads all policy rules from mongo collection.
func (a *adapter) LoadPolicy(model model.Model) (resErr error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objConfByID, err := a.findObjects(ctx)

	if err != nil {
		return err
	}

	roleNamesByID, err := a.loadRoles(ctx, model, objConfByID)

	if err != nil {
		return err
	}

	err = a.loadSubjects(ctx, model, roleNamesByID)

	if err != nil {
		return err
	}

	return nil
}

// SavePolicy isn't implemented.
// Implement it if all API is migrated to Go
// and there is not time to refactor mongo collection.
func (adapter) SavePolicy(model.Model) error {
	panic("implement me")
}

func (adapter) AddPolicy(string, string, []string) error {
	panic("implement me")
}

func (adapter) RemovePolicy(string, string, []string) error {
	panic("implement me")
}

func (adapter) RemoveFilteredPolicy(string, string, int, ...string) error {
	panic("implement me")
}

// findObjects fetches objects from mongo collection and returns map[objectID]objectConfig.
func (a *adapter) findObjects(ctx context.Context) (objConfByID map[string]objectConfig, resErr error) {
	cursor, err := a.collection.Find(
		ctx,
		bson.M{
			"crecord_type": libmodel.LineTypeObject,
			"crecord_name": bson.M{"$exists": true, "$ne": ""},
		},
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil && resErr == nil {
			resErr = err
		}
	}()

	objConfByID = make(map[string]objectConfig)

	for cursor.Next(ctx) {
		var line libmodel.Rbac
		err = cursor.Decode(&line)

		if err != nil {
			return nil, err
		}

		objConfByID[line.ID] = objectConfig{
			Name: line.Name,
			Type: line.ObjectType,
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return objConfByID, nil
}

// loadRoles fetches roles from mongo collection and adds them to casbin policy.
// Method returns map[roleID]roleName.
func (a *adapter) loadRoles(
	ctx context.Context,
	model model.Model,
	objConfByID map[string]objectConfig,
) (roleNamesByID map[string]string, resErr error) {
	cursor, err := a.collection.Find(
		ctx,
		bson.M{
			"crecord_type": libmodel.LineTypeRole,
			"crecord_name": bson.M{"$exists": true, "$ne": ""},
		},
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil && resErr == nil {
			resErr = err
		}
	}()

	roleNamesByID = make(map[string]string)
	ptype := CasbinPtypePolicy
	sec := ptype
	permBitmasksByType := map[string]map[string]int{
		libmodel.LineObjectTypeCRUD: {
			libmodel.PermissionCreate: libmodel.PermissionBitmaskCreate,
			libmodel.PermissionRead:   libmodel.PermissionBitmaskRead,
			libmodel.PermissionUpdate: libmodel.PermissionBitmaskUpdate,
			libmodel.PermissionDelete: libmodel.PermissionBitmaskDelete,
		},
		libmodel.LineObjectTypeRW: {
			libmodel.PermissionRead:   libmodel.PermissionBitmaskRead,
			libmodel.PermissionUpdate: libmodel.PermissionBitmaskUpdate,
			libmodel.PermissionDelete: libmodel.PermissionBitmaskDelete,
		},
	}

	for cursor.Next(ctx) {
		var line libmodel.Rbac
		err = cursor.Decode(&line)

		if err != nil {
			return nil, err
		}

		roleName := line.Name

		for objId, permConfig := range line.PermConfigList {
			if objConf, ok := objConfByID[objId]; ok {
				if objConf.Type != "" {
					if permBitmasksByName, ok := permBitmasksByType[objConf.Type]; ok {
						for permName, bitmask := range permBitmasksByName {
							if permConfig.Bitmask&bitmask == bitmask {
								model.AddPolicy(sec, ptype, []string{roleName, objConf.Name, permName})
							}
						}
					} else {
						return nil, fmt.Errorf("unknown config type \"%s\"", objConf.Type)
					}
				} else if permConfig.Bitmask&libmodel.PermissionBitmaskCan == libmodel.PermissionBitmaskCan {
					model.AddPolicy(sec, ptype, []string{roleName, objConf.Name, libmodel.PermissionCan})
				}
			}
		}

		roleNamesByID[line.ID] = roleName
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return roleNamesByID, nil
}

// loadSubjects loads subjects from mongo collection and adds them to casbin policy.
func (a *adapter) loadSubjects(
	ctx context.Context,
	model model.Model,
	roleNamesByID map[string]string,
) (resErr error) {
	cursor, err := a.collection.Find(ctx, bson.M{
		"crecord_type": libmodel.LineTypeSubject,
		"role":         bson.M{"$exists": true, "$ne": ""},
	})

	if err != nil {
		return err
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil && resErr == nil {
			resErr = err
		}
	}()

	ptype := CasbinPtypeRole
	sec := ptype

	for cursor.Next(ctx) {
		var line libmodel.Rbac
		err := cursor.Decode(&line)

		if err != nil {
			return err
		}

		subjectID := line.ID

		if roleName, ok := roleNamesByID[line.Role]; ok {
			model.AddPolicy(sec, ptype, []string{subjectID, roleName})
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}
