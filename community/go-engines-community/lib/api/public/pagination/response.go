package pagination

import "math"

type ListResponse[T any] struct {
	Data []T          `json:"data"`
	Meta MetaResponse `json:"meta"`
}

type MetaResponse struct {
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

type Data[T any] struct {
	Data       []T   `bson:"data" json:"data"`
	TotalCount int64 `bson:"total_count" json:"total_count"`
}

func NewResponse[T any](q Query, d Data[T]) ListResponse[T] {
	data := d.Data
	if data == nil {
		data = []T{}
	}

	return ListResponse[T]{
		Data: data,
		Meta: NewMeta(q, d.TotalCount),
	}
}

func NewMeta(q Query, total int64) MetaResponse {
	var pageCount int64
	if total > q.Limit && q.Limit > 0 {
		pageCount = int64(math.Ceil(float64(total) / float64(q.Limit)))
	} else {
		pageCount = 1
	}

	return MetaResponse{
		Page:       q.Page,
		PerPage:    q.Limit,
		PageCount:  pageCount,
		TotalCount: total,
	}
}
