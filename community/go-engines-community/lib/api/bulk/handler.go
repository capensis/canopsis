package bulk

import (
	"encoding/json"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
)

// Handler parses and validates a bulk request
// and calls f for each item in the request payload.
//
// Example:
//
// Suppose you receive a bulk request body like:
//
// [ {"_id": "A", "amount": 10}, {"_id": "B", "amount": 20}, {"_id": "C", "amount": 5} ]
//
// And you want to process each item individually:
//
//	type RequestItem struct {
//	    ID     string `json:"_id" binding:"required"`
//	    Amount int    `json:"amount" binding:"required,gt=0"`
//	}
//
//	f := func(i RequestItem) (string, error) {
//	    // Do something with each item.
//	    // Return ID and error.
//	    err := saveItemToDB(i)
//
//	    return i.ID, err
//	}
func Handler[T any](
	c *gin.Context,
	f func(T) (string, error),
	responder httperror.Responder,
) {
	raw, err := c.GetRawData()
	if err != nil {
		responder.Respond(c, err)

		return
	}

	parsed, err := fastjson.ParseBytes(raw)
	if err != nil {
		responder.Respond(c, validation.ErrInvalidRequestBody)

		return
	}

	items, err := parsed.Array()
	if err != nil {
		responder.Respond(c, validation.ErrInvalidRequestBody)

		return
	}

	var ar fastjson.Arena
	res := ar.NewArray()
	for idx, item := range items {
		req, err := validateItem[T](item)
		if err != nil {
			res.SetArrayItem(idx, getErrResponseItem(c, err, item, responder, &ar))
			continue
		}

		id, err := f(req)
		if err != nil {
			res.SetArrayItem(idx, getErrResponseItem(c, err, item, responder, &ar))
			continue
		}

		res.SetArrayItem(idx, getResponseItem(http.StatusOK, item, id, &ar))
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, res.MarshalTo(nil))
}

// HandlerWithGrouping parses and validates a bulk request,
// groups items using mergeReqs after comparing them with compareWithPrev,
// and then calls procReq for each resulting group.
// Use this to reduce the number of database queries.
//
// Example:
//
// Suppose you receive a bulk request of items that look like this:
//
// [ {"_id": "A", "amount": 10}, {"_id": "A", "amount": 20}, {"_id": "B", "amount": 5} ]
//
// You can group consecutive items with the same _id and sum the amounts:
//
//	compareWithPrev := func(prev, cur RequestItem) bool {
//	    return prev.ID == cur.ID
//	}
//
//	mergeReqs := func(merged, cur RequestItem) RequestItem {
//	    merged.Amount += cur.Amount
//
//	    return merged
//	}
//
//	procReq := func(i RequestItem) (string, error) {
//	    // Process each merged group (A:30 and B:5)
//	    err := saveToDB(i)
//
//	    return i.ID, err
//	}
func HandlerWithGrouping[T any](
	c *gin.Context,
	compareWithPrev func(prev, cur T) bool,
	mergeReqs func(merged, cur T) T,
	procReq func(T) (id string, err error),
	responder httperror.Responder,
) {
	raw, err := c.GetRawData()
	if err != nil {
		responder.Respond(c, err)

		return
	}

	parsed, err := fastjson.ParseBytes(raw)
	if err != nil {
		responder.Respond(c, validation.ErrInvalidRequestBody)

		return
	}

	items, err := parsed.Array()
	if err != nil {
		responder.Respond(c, validation.ErrInvalidRequestBody)

		return
	}

	var ar fastjson.Arena
	res := ar.NewArray()
	type groupedRequest struct {
		Req   T
		Items map[int]*fastjson.Value
	}

	groupedReqs := make([]groupedRequest, 0)
	var prevReq T
	for idx, item := range items {
		req, err := validateItem[T](item)
		if err != nil {
			res.SetArrayItem(idx, getErrResponseItem(c, err, item, responder, &ar))
			continue
		}

		if idx > 0 && compareWithPrev(prevReq, req) {
			groupedReqs[len(groupedReqs)-1].Req = mergeReqs(groupedReqs[len(groupedReqs)-1].Req, req)
			groupedReqs[len(groupedReqs)-1].Items[idx] = item
		} else {
			groupedReqs = append(groupedReqs, groupedRequest{Req: req, Items: map[int]*fastjson.Value{idx: item}})
		}

		prevReq = req
	}

	for _, g := range groupedReqs {
		req := g.Req
		id, err := procReq(req)
		if err != nil {
			for idx, item := range g.Items {
				res.SetArrayItem(idx, getErrResponseItem(c, err, item, responder, &ar))
			}

			continue
		}

		for idx, item := range g.Items {
			res.SetArrayItem(idx, getResponseItem(http.StatusOK, item, id, &ar))
		}
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, res.MarshalTo(nil))
}

func validateItem[T any](item *fastjson.Value) (T, error) {
	var r T
	obj, err := item.Object()
	if err != nil {
		return r, validation.ErrInvalidRequestBody
	}

	err = json.Unmarshal(obj.MarshalTo(nil), &r)
	if err != nil {
		return r, validation.ErrInvalidRequestBody
	}

	err = validation.ValidateStruct(r)

	return r, err
}

func getResponseItem(statusCode int, item *fastjson.Value, id string, ar *fastjson.Arena) *fastjson.Value {
	res := ar.NewObject()
	res.Set("status", ar.NewNumberInt(statusCode))
	res.Set("item", item)
	res.Set("id", ar.NewString(id))

	return res
}

func getErrResponseItem(c *gin.Context, err error, item *fastjson.Value, responder httperror.Responder, ar *fastjson.Arena) *fastjson.Value {
	code, res := responder.GetResponse(c, err)
	res.Set("status", ar.NewNumberInt(code))
	res.Set("item", item)

	return res
}
