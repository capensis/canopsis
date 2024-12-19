package role

import (
	"context"
	"errors"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
	ValidateBulkUpdatePermissionsRequestItem(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	dbCollection           mongo.DbCollection
	dbPermissionCollection mongo.DbCollection
	dbViewCollection       mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbCollection:           dbClient.Collection(mongo.RoleCollection),
		dbPermissionCollection: dbClient.Collection(mongo.PermissionCollection),
		dbViewCollection:       dbClient.Collection(mongo.ViewMongoCollection),
	}
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	// Validate default view
	if r.DefaultView != "" {
		err := v.dbViewCollection.FindOne(ctx, bson.M{"_id": r.DefaultView}).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				sl.ReportError(r.DefaultView, "DefaultView", "DefaultView", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}

	v.validatePermissions(ctx, sl, r.Type, r.Permissions)
}

func (v *baseValidator) ValidateBulkUpdatePermissionsRequestItem(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkUpdatePermissionsRequestItem)
	var roleType string
	if r.ID != "" {
		var role Response
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}, options.FindOne().SetProjection(bson.M{"type": 1})).
			Decode(&role)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			panic(err)
		}

		roleType = role.Type
	}

	v.validatePermissions(ctx, sl, roleType, r.Permissions)
}

func (v *baseValidator) validatePermissions(ctx context.Context, sl validator.StructLevel, roleType string, permissions map[string][]string) {
	if len(permissions) == 0 {
		return
	}

	perms, _, err := getPermissions(ctx, v.dbPermissionCollection, permissions)
	if err != nil {
		panic(err)
	}

	for id, actions := range permissions {
		switch roleType {
		case TypeUI:
			if strings.HasPrefix(id, PermissionAPIPrefix) {
				sl.ReportError(permissions[id], "Permissions."+id, "Permissions."+id, "not_ui_perm", "")
			}
		case TypeAPI:
			if !strings.HasPrefix(id, PermissionAPIPrefix) {
				sl.ReportError(permissions[id], "Permissions."+id, "Permissions."+id, "not_api_perm", "")
			}
		}

		if perm, ok := perms[id]; ok {
			var validActions []string
			switch perm.Type {
			case "":
				if len(actions) > 0 {
					sl.ReportError(permissions[id], "Permissions."+id, "Permissions."+id, "must_be_empty", "")
				}
			case securitymodel.ObjectTypeCRUD:
				validActions = []string{
					securitymodel.PermissionCreate,
					securitymodel.PermissionRead,
					securitymodel.PermissionUpdate,
					securitymodel.PermissionDelete,
				}
			case securitymodel.ObjectTypeRW:
				validActions = []string{
					securitymodel.PermissionRead,
					securitymodel.PermissionUpdate,
					securitymodel.PermissionDelete,
				}
			}

			if len(validActions) > 0 {
				for _, action := range actions {
					found := false
					for _, v := range validActions {
						if action == v {
							found = true
							break
						}
					}

					if !found {
						sl.ReportError(permissions[id], "Permissions."+id, "Permissions."+id, "oneof", strings.Join(validActions, " "))
					}
				}
			}
		} else {
			sl.ReportError(permissions[id], "Permissions."+id, "Permissions."+id, "not_exist", "")
		}
	}
}
