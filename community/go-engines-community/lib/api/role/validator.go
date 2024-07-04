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
)

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
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
	// Validate permissions
	if len(r.Permissions) == 0 {
		return
	}

	types, err := getTypes(ctx, v.dbPermissionCollection, r.Permissions)
	if err != nil {
		panic(err)
	}

	for id, actions := range r.Permissions {
		if t, ok := types[id]; ok {
			var validActions []string
			switch t {
			case "":
				if len(actions) > 0 {
					sl.ReportError(r.Permissions[id], "Permissions."+id, "Permissions."+id, "must_be_empty", "")
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
						sl.ReportError(r.Permissions[id], "Permissions."+id, "Permissions."+id, "oneof", strings.Join(validActions, " "))
					}
				}
			}
		} else {
			sl.ReportError(r.Permissions[id], "Permissions."+id, "Permissions."+id, "not_exist", "")
		}
	}
}
