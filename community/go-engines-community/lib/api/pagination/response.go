package pagination

import "math"

type ListResponse struct {
	Data any          `json:"data"`
	Meta MetaResponse `json:"meta"`
}

type MetaResponse struct {
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

type Data interface {
	GetData() interface{}
	GetTotal() int64
}

func NewResponse(q Query, d Data) ListResponse {
	data := d.GetData()
	if data == nil {
		data = []any{}
	}

	return ListResponse{
		Data: data,
		Meta: NewMeta(q, d.GetTotal()),
	}
}

func NewMeta(q Query, total int64) MetaResponse {
	if !q.Paginate {
		q.Limit = total
	}

	var pageCount int64
	if q.Limit > 0 {
		pageCount = int64(math.Ceil(float64(total) / float64(q.Limit)))
	}

	if pageCount == 0 {
		pageCount = 1
	}

	return MetaResponse{
		Page:       q.Page,
		PerPage:    q.Limit,
		PageCount:  pageCount,
		TotalCount: total,
	}
}
