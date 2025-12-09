package pbehavior

import (
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.BulkAPI
	Patch(c *gin.Context)
	DeleteByName(c *gin.Context)
	ListByEntityID(c *gin.Context)
	CalendarByEntityID(c *gin.Context)
	ListEntities(c *gin.Context)
	BulkEntityCreate(c *gin.Context)
	BulkEntityDelete(c *gin.Context)
	BulkConnectorCreate(c *gin.Context)
	BulkConnectorDelete(c *gin.Context)
	BulkConnectorEdit(c *gin.Context)
	DBExport(c *gin.Context)
	ExecPattern(c *gin.Context)
	ExecAllPatterns(c *gin.Context)
}

type api struct {
	store          Store
	mongoExporter  dbexport.Exporter
	computeChan    chan<- rpc.PbehaviorRecomputeEvent
	jobPublisher   workers.JobPublisher
	errorResponder httperror.Responder
}

func NewApi(
	store Store,
	mongoExporter dbexport.Exporter,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	jobPublisher workers.JobPublisher,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:          store,
		mongoExporter:  mongoExporter,
		computeChan:    computeChan,
		jobPublisher:   jobPublisher,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// ListByEntityID
// @Success 200 {array} Response
func (a *api) ListByEntityID(c *gin.Context) {
	var r FindByEntityIDRequest
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	entity, err := a.store.FindEntity(c, r.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if entity == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	res, err := a.store.FindByEntityID(c, *entity, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}

// CalendarByEntityID
// @Success 200 {array} CalendarResponse
func (a *api) CalendarByEntityID(c *gin.Context) {
	var r CalendarByEntityIDRequest
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	entity, err := a.store.FindEntity(c, r.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if entity == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	res, err := a.store.CalendarByEntityID(c, *entity, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	pbh, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if pbh == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, pbh)
}

// ListEntities
// @Success 200 {object} pagination.ListResponse{data=[]entity.Entity}
func (a *api) ListEntities(c *gin.Context) {
	var r EntitiesListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.FindEntities(c, c.Param("id"), r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if aggregationResult == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	pbh, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{pbh.ID}, RecomputeInherited: pbh.Inherited})

	c.JSON(http.StatusCreated, pbh)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	pbh, recomputeInherited, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if pbh == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{pbh.ID}, RecomputeInherited: recomputeInherited})

	c.JSON(http.StatusOK, pbh)
}

// Patch
// @Param body body PatchRequest true "body"
// @Success 200 {object} Response
func (a *api) Patch(c *gin.Context) {
	request := PatchRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	pbh, recomputeInherited, err := a.store.UpdateByPatch(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if pbh == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{pbh.ID}, RecomputeInherited: recomputeInherited})

	c.JSON(http.StatusOK, pbh)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, recomputeInherited, err := a.store.Delete(c, id, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{id}, RecomputeInherited: recomputeInherited})
	c.JSON(http.StatusNoContent, nil)
}

func (a *api) DeleteByName(c *gin.Context) {
	request := DeleteByNameRequest{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	id, recomputeInherited, err := a.store.DeleteByName(c, request.Name, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if id == "" {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: []string{id}, RecomputeInherited: recomputeInherited})
	c.JSON(http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	ids := make([]string, 0)
	recomputeInherited := false

	bulk.Handler(c, func(request CreateRequest) (string, error) {
		pbh, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		ids = append(ids, pbh.ID)

		recomputeInherited = recomputeInherited || pbh.Inherited

		return pbh.ID, nil
	}, a.errorResponder)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: ids, RecomputeInherited: recomputeInherited})
	}
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	ids := make([]string, 0)
	exists := make(map[string]struct{})
	recomputeInherited := false

	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		pbh, curRecomputeInherited, err := a.store.Update(c, UpdateRequest(request))
		if err != nil {
			return "", err
		}

		if pbh == nil {
			return "", httperror.ErrNotFound
		}

		if _, ok := exists[pbh.ID]; !ok {
			ids = append(ids, pbh.ID)
			exists[pbh.ID] = struct{}{}
		}

		recomputeInherited = recomputeInherited || curRecomputeInherited

		return pbh.ID, nil
	}, a.errorResponder)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: ids, RecomputeInherited: recomputeInherited})
	}
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	recomputeInherited := false

	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, curRecomputeInherited, err := a.store.Delete(c, request.ID, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		ids = append(ids, request.ID)

		recomputeInherited = recomputeInherited || curRecomputeInherited

		return request.ID, nil
	}, a.errorResponder)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{Ids: ids, RecomputeInherited: recomputeInherited})
	}
}

// BulkEntityCreate
// @Param body body []BulkEntityCreateRequestItem true "body"
func (a *api) BulkEntityCreate(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkEntityCreateRequestItem) (string, error) {
		pbh, err := a.store.EntityInsert(c, request)
		if err != nil {
			return "", err
		}

		if pbh == nil {
			return "", httperror.ErrNotFound
		}

		ids = append(ids, pbh.ID)

		return pbh.ID, nil
	}, a.errorResponder)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    username,
			UserID:    userID,
			Initiator: types.InitiatorUser,
		})
	}
}

// BulkEntityDelete
// @Param body body []BulkEntityDeleteRequestItem true "body"
func (a *api) BulkEntityDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	ids := make([]string, 0)
	bulk.Handler(c, func(request BulkEntityDeleteRequestItem) (string, error) {
		id, err := a.store.EntityDelete(c, request)
		if err != nil {
			return "", err
		}

		if id == "" {
			return "", httperror.ErrNotFound
		}

		ids = append(ids, id)

		return id, nil
	}, a.errorResponder)

	if len(ids) > 0 {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    username,
			UserID:    userID,
			Initiator: types.InitiatorUser,
		})
	}
}

// BulkConnectorCreate
// @Param body body []BulkConnectorCreateRequestItem true "body"
func (a *api) BulkConnectorCreate(c *gin.Context) {
	idsByOrigin := make(map[string][]string)
	exists := make(map[string]struct{})
	bulk.HandlerWithGrouping(c,
		func(prev, cur BulkConnectorCreateRequestItem) bool {
			return cur.Origin == prev.Origin &&
				cur.Comment == prev.Comment &&
				cur.Start.Unix() == prev.Start.Unix() &&
				cur.Stop.Unix() == prev.Stop.Unix()
		},
		func(merged, cur BulkConnectorCreateRequestItem) BulkConnectorCreateRequestItem {
			merged.Entities = append(merged.Entities, cur.Entities...)

			return merged
		},
		func(request BulkConnectorCreateRequestItem) (string, error) {
			pbh, err := a.store.ConnectorCreate(c, request)
			if err != nil {
				return "", err
			}

			if pbh == nil {
				return "", httperror.ErrNotFound
			}

			if _, ok := exists[pbh.ID]; !ok {
				// skip add to idsByOrigin when pbh is not effective
				if pbh.Stop == nil || pbh.RRule != "" || pbh.Stop.Unix() >= time.Now().Unix() {
					idsByOrigin[request.Origin] = append(idsByOrigin[request.Origin], pbh.ID)
				}
				exists[pbh.ID] = struct{}{}
			}

			return pbh.ID, nil
		},
		a.errorResponder,
	)

	for origin, ids := range idsByOrigin {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

// BulkConnectorDelete
// @Param body body []BulkConnectorDeleteRequestItem true "body"
func (a *api) BulkConnectorDelete(c *gin.Context) {
	idsByOrigin := make(map[string][]string)
	exists := make(map[string]struct{})
	bulk.HandlerWithGrouping(c,
		func(prev, cur BulkConnectorDeleteRequestItem) bool {
			return cur.Origin == prev.Origin &&
				cur.Comment == prev.Comment &&
				cur.Start.Unix() == prev.Start.Unix() &&
				cur.Stop.Unix() == prev.Stop.Unix()
		},
		func(merged, cur BulkConnectorDeleteRequestItem) BulkConnectorDeleteRequestItem {
			merged.Entities = append(merged.Entities, cur.Entities...)

			return merged
		},
		func(request BulkConnectorDeleteRequestItem) (string, error) {
			id, err := a.store.ConnectorDelete(c, request)
			if err != nil {
				return "", err
			}

			if id == "" {
				return "", httperror.ErrNotFound
			}

			if _, ok := exists[id]; !ok {
				idsByOrigin[request.Origin] = append(idsByOrigin[request.Origin], id)
				exists[id] = struct{}{}
			}

			return id, nil
		},
		a.errorResponder,
	)

	for origin, ids := range idsByOrigin {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

// BulkConnectorEdit
// @Param body body []BulkConnectorEditRequestItem true "body"
func (a *api) BulkConnectorEdit(c *gin.Context) {
	idsByOrigin := make(map[string][]string)
	exists := make(map[string]struct{})
	bulk.HandlerWithGrouping(c,
		func(prev, cur BulkConnectorEditRequestItem) bool {
			return cur.Action == prev.Action &&
				cur.Origin == prev.Origin &&
				cur.Comment == prev.Comment &&
				cur.Start.Unix() == prev.Start.Unix() &&
				cur.Stop.Unix() == prev.Stop.Unix()
		},
		func(merged, cur BulkConnectorEditRequestItem) BulkConnectorEditRequestItem {
			merged.Entities = append(merged.Entities, cur.Entities...)

			return merged
		},
		func(request BulkConnectorEditRequestItem) (string, error) {
			var id string
			var err error
			switch request.Action {
			case BulkConnectorActionCreate:
				var pbh *Response
				pbh, err = a.store.ConnectorCreate(c, BulkConnectorCreateRequestItem{
					Author:   request.Author,
					Entities: request.Entities,
					Origin:   request.Origin,
					Start:    request.Start,
					Stop:     request.Stop,
					Comment:  request.Comment,
					Name:     request.Name,
					Reason:   request.Reason,
					Type:     request.Type,
					Color:    request.Color,
				})
				if pbh != nil {
					id = pbh.ID
				}
			case BulkConnectorActionDelete:
				id, err = a.store.ConnectorDelete(c, BulkConnectorDeleteRequestItem{
					Author:   request.Author,
					Entities: request.Entities,
					Origin:   request.Origin,
					Start:    request.Start,
					Stop:     request.Stop,
					Comment:  request.Comment,
				})
			default:
				return "", validation.NewSingleErrorWithParam("oneof", "Action", "Action", BulkConnectorActionCreate+" "+BulkConnectorActionDelete, request)
			}

			if err != nil {
				return "", err
			}

			if id == "" {
				return "", httperror.ErrNotFound
			}

			if _, ok := exists[id]; !ok {
				idsByOrigin[request.Origin] = append(idsByOrigin[request.Origin], id)
				exists[id] = struct{}{}
			}

			return id, nil
		},
		a.errorResponder,
	)

	for origin, ids := range idsByOrigin {
		a.sendComputeTask(rpc.PbehaviorRecomputeEvent{
			Ids:       ids,
			Author:    origin,
			Initiator: types.InitiatorExternal,
		})
	}
}

// DBExport
// @Param body body dbexport.Request true "body"
func (a *api) DBExport(c *gin.Context) {
	request := dbexport.Request{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	b, err := a.mongoExporter.Export(c, mongo.PbehaviorMongoCollection, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	dbexport.AttachFile(c, mongo.PbehaviorMongoCollection, b)
}

// ExecPattern
// @Param body body ExecPatternRequest true "body"
// @Success 200 {object} pattern.CountResponse
func (a *api) ExecPattern(c *gin.Context) {
	request := ExecPatternRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.ExecPatternAndUpdate(c, request.ID, request.EntityPattern)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) ExecAllPatterns(c *gin.Context) {
	err := a.jobPublisher.Publish(c, "")
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) sendComputeTask(event rpc.PbehaviorRecomputeEvent) {
	a.computeChan <- event
}
