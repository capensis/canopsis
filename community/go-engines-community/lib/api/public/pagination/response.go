package pagination

type ListResponse[T any] struct {
	Data []T          `json:"data"`
	Meta MetaResponse `json:"meta"`
}

type MetaResponse struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

func NewResponse[T any](q Query, d []T) ListResponse[T] {
	data := d
	if data == nil {
		data = []T{}
	}

	if len(data) > q.Limit {
		data = data[:q.Limit]
	}

	return ListResponse[T]{
		Data: data,
		Meta: NewMeta(q, d),
	}
}

func NewMeta[T any](q Query, d []T) MetaResponse {
	return MetaResponse{
		Limit:   q.Limit,
		Offset:  q.Offset,
		HasMore: len(d) > q.Limit,
	}
}
