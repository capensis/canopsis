package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const apiRoute = "/api/v2/links"

func NewGenerator(
	legacyUrl string,
	dbClient mongo.DbClient,
	httpClient libhttp.Doer,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) liblink.Generator {
	return &generator{
		legacyUrl:        legacyUrl,
		httpClient:       httpClient,
		encoder:          encoder,
		decoder:          decoder,
		alarmCollection:  dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection: dbClient.Collection(mongo.EntityMongoCollection),
	}
}

type generator struct {
	legacyUrl        string
	httpClient       libhttp.Doer
	encoder          encoding.Encoder
	decoder          encoding.Decoder
	alarmCollection  mongo.DbCollection
	entityCollection mongo.DbCollection
}

type fetchLinksRequest struct {
	Entities []fetchLinksRequestItem `json:"entities"`
}

type fetchLinksRequestItem struct {
	Alarm  string `json:"alarm"`
	Entity string `json:"entity"`
}

type fetchLinksResponse struct {
	Data []fetchLinksResponseItem `json:"data"`
}

type fetchLinksResponseItem struct {
	fetchLinksRequestItem
	Links map[string][]fetchLinksResponseLink `json:"links"`
}

type fetchLinksResponseLink struct {
	Label string `json:"label"`
	Link  string `json:"link"`
}

func (g *generator) Load(_ context.Context) error {
	return nil
}

func (g *generator) GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user liblink.User) (liblink.LinksByCategory, error) {
	res, err := g.GenerateForAlarms(ctx, []string{alarm.ID}, user)
	if err != nil {
		return nil, err
	}

	return res[alarm.ID], nil
}

func (g *generator) GenerateForAlarms(ctx context.Context, ids []string, _ liblink.User) (map[string]liblink.LinksByCategory, error) {
	req, err := g.createRequestByAlarms(ctx, ids)
	if err != nil || req == nil {
		return nil, err
	}

	data, err := g.doRequest(req)
	if err != nil {
		return nil, err
	}

	res := make(map[string]liblink.LinksByCategory, len(data))
	for _, v := range data {
		res[v.Alarm] = make(liblink.LinksByCategory, len(v.Links))
		for category, links := range v.Links {
			res[v.Alarm][category] = g.transformLinks(links)
		}
	}

	return res, nil
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string, _ liblink.User) (map[string]liblink.LinksByCategory, error) {
	req, err := g.createRequestByEntities(ctx, ids)
	if err != nil || req == nil {
		return nil, err
	}

	data, err := g.doRequest(req)
	if err != nil {
		return nil, err
	}

	res := make(map[string]liblink.LinksByCategory, len(data))
	for _, v := range data {
		res[v.Entity] = make(liblink.LinksByCategory, len(v.Links))
		for category, links := range v.Links {
			res[v.Entity][category] = g.transformLinks(links)
		}
	}

	return res, nil
}

func (g *generator) GenerateCombinedForAlarmsByRule(_ context.Context, _ string, _ []string, _ liblink.User) ([]liblink.Link, error) {
	return nil, nil
}

func (g *generator) createRequestByAlarms(ctx context.Context, ids []string) (*http.Request, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	cursor, err := g.alarmCollection.Find(ctx, bson.M{
		"_id":        bson.M{"$in": ids},
		"v.resolved": nil,
	}, options.Find().SetProjection(bson.M{"d": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]fetchLinksRequestItem, 0, len(ids))
	for cursor.Next(ctx) {
		alarm := types.Alarm{}
		err = cursor.Decode(&alarm)
		if err != nil {
			return nil, err
		}
		items = append(items, fetchLinksRequestItem{
			Alarm:  alarm.ID,
			Entity: alarm.EntityID,
		})
	}

	return g.createRequest(ctx, items)
}

func (g *generator) createRequestByEntities(ctx context.Context, ids []string) (*http.Request, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	cursor, err := g.entityCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"pipeline": []bson.M{
				{"$match": bson.M{"v.resolved": nil}},
				{"$project": bson.M{
					"_id": 1,
				}},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$project": bson.M{
			"alarm":  "$alarm._id",
			"entity": "$_id",
		}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]fetchLinksRequestItem, 0, len(ids))
	for cursor.Next(ctx) {
		res := struct {
			AlarmID  string `bson:"alarm"`
			EntityID string `bson:"entity"`
		}{}
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		items = append(items, fetchLinksRequestItem{
			Alarm:  res.AlarmID,
			Entity: res.EntityID,
		})
	}

	return g.createRequest(ctx, items)
}

func (g *generator) createRequest(ctx context.Context, items []fetchLinksRequestItem) (*http.Request, error) {
	if len(items) == 0 {
		return nil, nil
	}

	body, err := g.encoder.Encode(fetchLinksRequest{Entities: items})
	if err != nil {
		return nil, err
	}

	u, err := url.JoinPath(g.legacyUrl, apiRoute)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", canopsis.JsonContentType)
	if c, ok := ctx.(*gin.Context); ok {
		// Add old v3 API auth credentials
		req.Header.Add(security.HeaderApiKey, c.GetString(auth.ApiKey))
	}
	return req, nil
}

func (g *generator) doRequest(req *http.Request) ([]fetchLinksResponseItem, error) {
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code=%d body=%q", resp.StatusCode, string(buf))
	}

	var resBody fetchLinksResponse
	err = g.decoder.Decode(buf, &resBody)
	if err != nil {
		return nil, err
	}

	return resBody.Data, nil
}

func (g *generator) transformLinks(responseLinks []fetchLinksResponseLink) []liblink.Link {
	links := make([]liblink.Link, len(responseLinks))
	for i, l := range responseLinks {
		links[i] = liblink.Link{
			Label:  l.Label,
			Url:    l.Link,
			Action: liblink.ActionOpen,
		}
	}
	return links
}
