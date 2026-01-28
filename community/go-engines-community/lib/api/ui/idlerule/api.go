package idlerule

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.BulkAPI
	DBExport(c *gin.Context)
}

type api struct {
	store          Store
	mongoExporter  dbexport.Exporter
	errorResponder response.ErrorResponder
}

func NewApi(
	store Store,
	mongoExporter dbexport.Exporter,
	errorResponder response.ErrorResponder,
) API {
	return &api{
		store:          store,
		mongoExporter:  mongoExporter,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]idlerule.Rule}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := binding.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rules, err := a.store.Find(c, query)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, rules)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} idlerule.Rule
func (a *api) Get(c *gin.Context) {
	rule, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if rule == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, rule)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} idlerule.Rule
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rule, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if rule == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusCreated, rule)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Rule
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	rule, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if rule == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, rule)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.Status(http.StatusNoContent)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		rule, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		return rule.ID, nil
	}, a.errorResponder)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		rule, err := a.store.Update(c, UpdateRequest(request))
		if err != nil {
			return "", err
		}

		if rule == nil {
			return "", httperror.ErrNotFound
		}

		return rule.ID, nil
	}, a.errorResponder)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		return request.ID, nil
	}, a.errorResponder)
}

// DBExport
// @Param body body dbexport.Request true "body"
func (a *api) DBExport(c *gin.Context) {
	request := dbexport.Request{}

	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	b, err := a.mongoExporter.Export(c, mongo.IdleRuleMongoCollection, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	dbexport.AttachFile(c, mongo.IdleRuleMongoCollection, b)
}
