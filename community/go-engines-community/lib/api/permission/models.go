package permission

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=name description"`
}

type Permission struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Type        string `bson:"type" json:"type"`
	// View is used only for view permissions.
	View *View `bson:"view" json:"view,omitempty"`
	// ViewGroup is used only for view permissions.
	ViewGroup *ViewGroup `bson:"view_group" json:"view_group,omitempty"`
	// Playlist is used only for playlist permissions.
	Playlist *Playlist `bson:"playlist" json:"playlist,omitempty"`
}

type View struct {
	ID       string `bson:"_id" json:"_id"`
	Title    string `bson:"title" json:"title"`
	Position int64  `bson:"position" json:"position"`
}

type ViewGroup struct {
	ID       string `bson:"_id" json:"_id"`
	Title    string `bson:"title" json:"title"`
	Position int64  `bson:"position" json:"position"`
}

type Playlist struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
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
