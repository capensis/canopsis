package permission

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"

type ListRequest struct {
	pagination.FilteredQuery
}

type Permission struct {
	ID          string  `bson:"_id" json:"_id"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	Type        string  `bson:"type" json:"type"`
	Groups      []Group `bson:"groups" json:"groups"`
}

type Group struct {
	ID       string `bson:"_id" json:"_id"`
	Title    string `bson:"title" json:"title,omitempty"`
	Position int64  `bson:"position" json:"position"`
}

type AggregationResult struct {
	Data       []Permission `bson:"data" json:"data"`
	TotalCount int64        `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
