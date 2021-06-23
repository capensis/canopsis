package heartbeat

type ModelTransformer interface {
	TransformCreateRequestToModel(request CreateRequest) *Response
	TransformUpdateRequestToModel(request UpdateRequest) *Response
	TransformBulkCreateRequestToModels(request BulkCreateRequest) []*Response
	TransformBulkUpdateRequestToModels(request BulkUpdateRequest) []*Response
}

func NewModelTransformer() ModelTransformer {
	return &modelTransformer{}
}

type modelTransformer struct{}

func (modelTransformer) TransformCreateRequestToModel(request CreateRequest) *Response {
	return &Response{
		ID:               request.ID,
		Name:             request.Name,
		Description:      request.Description,
		Author:           request.Author,
		Pattern:          request.Pattern,
		ExpectedInterval: request.ExpectedInterval,
		Output:           request.Output,
	}
}

func (t *modelTransformer) TransformUpdateRequestToModel(request UpdateRequest) *Response {
	return t.TransformCreateRequestToModel(CreateRequest{
		BaseEditRequest: request.BaseEditRequest,
		ID:              request.ID,
	})
}

func (t *modelTransformer) TransformBulkCreateRequestToModels(request BulkCreateRequest) []*Response {
	heartbeats := make([]*Response, len(request.Items))
	for i := range heartbeats {
		heartbeats[i] = t.TransformCreateRequestToModel(request.Items[i])
	}

	return heartbeats
}

func (t *modelTransformer) TransformBulkUpdateRequestToModels(request BulkUpdateRequest) []*Response {
	heartbeats := make([]*Response, len(request.Items))
	for i := range heartbeats {
		heartbeats[i] = t.TransformCreateRequestToModel(CreateRequest{
			BaseEditRequest: request.Items[i].BaseEditRequest,
			ID:              request.Items[i].ID,
		})
	}

	return heartbeats
}
