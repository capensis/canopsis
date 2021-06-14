package role

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type Validator interface {
	ValidateCreateRequest(ctx context.Context, sl validator.StructLevel)
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.RightsMongoCollection),
	}
}

func (v *baseValidator) ValidateCreateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	// Validate name
	if r.Name != "" {
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.Name}).Err()
		if err == nil {
			sl.ReportError(r.Name, "Name", "Name", "unique", "")
		} else if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	// Validate default view
	if r.DefaultView != "" {
		err := v.dbClient.Collection(mongo.ViewMongoCollection).FindOne(ctx, bson.M{"_id": r.DefaultView}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(r.DefaultView, "DefaultView", "DefaultView", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
	// Validate permissions
	if len(r.Permissions) != 0 {
		ids := make([]string, 0)
		for id := range r.Permissions {
			ids = append(ids, id)
		}
		cursor, err := v.dbCollection.Find(ctx, bson.M{
			"_id":          bson.M{"$in": ids},
			"crecord_type": securitymodel.LineTypeObject,
		})
		if err != nil {
			panic(err)
		}

		defer cursor.Close(ctx)

		res := make([]struct {
			ID   string `bson:"_id"`
			Type string `bson:"type"`
		}, 0)
		err = cursor.All(ctx, &res)
		if err != nil {
			panic(err)
		}

		types := make(map[string]string)
		for _, v := range res {
			types[v.ID] = v.Type
		}

		for id, actions := range r.Permissions {
			if t, ok := types[id]; ok {
				var validActions []string
				switch t {
				case "":
					if len(actions) > 0 {
						sl.ReportError(r.Permissions[id], "Permissions."+id, "Permissions."+id, "must_be_empty", "")
					}
				case securitymodel.LineObjectTypeCRUD:
					validActions = []string{
						securitymodel.PermissionCreate,
						securitymodel.PermissionRead,
						securitymodel.PermissionUpdate,
						securitymodel.PermissionDelete,
					}
				case securitymodel.LineObjectTypeRW:
					validActions = []string{
						securitymodel.PermissionRead,
						securitymodel.PermissionUpdate,
						securitymodel.PermissionDelete,
					}
				}

				if len(validActions) > 0 {
					if len(actions) == 0 {
						sl.ReportError(r.Permissions[id], "Permissions."+id, "Permissions."+id, "oneof", strings.Join(validActions, " "))
					} else {
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
				}
			} else {
				sl.ReportError(r.Permissions[id], "Permissions."+id, "Permissions."+id, "not_exist", "")
			}
		}
	}
}
