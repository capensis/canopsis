package view

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/junit"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
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

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	// Validate group
	if r.Group != "" {
		err := v.dbClient.Collection(mongo.ViewGroupMongoCollection).FindOne(ctx, bson.M{"_id": r.Group}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(r.Group, "Group", "Group", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}

func ValidateWidgetParametersJunitRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(WidgetParametersJunitRequest)
	isAPI := r.IsAPI

	if r.Directory == "" {
		if !isAPI {
			sl.ReportError(r.Directory, "Directory", "Directory", "required", "")
		}
	} else if isAPI {
		sl.ReportError(r.Directory, "Directory", "Directory", "must_be_empty", "")
	}

	if len(r.ScreenshotDirectories) > 0 && isAPI {
		sl.ReportError(r.ScreenshotDirectories, "ScreenshotDirectories", "ScreenshotDirectories", "must_be_empty", "")
	}

	if len(r.VideoDirectories) > 0 && isAPI {
		sl.ReportError(r.VideoDirectories, "VideoDirectories", "VideoDirectories", "must_be_empty", "")
	}

	if r.ScreenshotFilemask != "" {
		_, err := junit.NewFileMask(r.ScreenshotFilemask)
		if err != nil {
			sl.ReportError(r.ScreenshotFilemask, "ScreenshotFilemask", "ScreenshotFilemask", "filemask", "")
		}
	}

	if r.VideoFilemask != "" {
		_, err := junit.NewFileMask(r.VideoFilemask)
		if err != nil {
			sl.ReportError(r.VideoFilemask, "VideoFilemask", "VideoFilemask", "filemask", "")
		}
	}
}

func ValidateEditPositionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditPositionRequest)

	if len(r.Items) > 0 {
		exists := make(map[string]bool, len(r.Items))
		existsView := make(map[string]bool, len(r.Items))
		for _, item := range r.Items {
			if exists[item.ID] {
				sl.ReportError(r.Items, "Items", "Item", "has_duplicates", "")
				return
			}

			exists[item.ID] = true

			for _, view := range item.Views {
				if existsView[view] {
					sl.ReportError(r.Items, "Items", "Item", "has_duplicates", "")
					return
				}

				existsView[view] = true
			}
		}
	}
}
