package heartbeat

type ModelTransformer interface {
	TransformCreateRequestToModel(request CreateRequest) *Heartbeat
	TransformUpdateRequestToModel(request UpdateRequest) *Heartbeat
	TransformBulkCreateRequestToModels(request BulkCreateRequest) []*Heartbeat
	TransformBulkUpdateRequestToModels(request BulkUpdateRequest) []*Heartbeat
}

func NewModelTransformer() ModelTransformer {
	return &modelTransformer{}
}

type modelTransformer struct{}

func (modelTransformer) TransformCreateRequestToModel(request CreateRequest) *Heartbeat {
	return &Heartbeat{
		ID:				  request.ID,
		Name:             request.Name,
		Description:      request.Description,
		Author:           request.Author,
		Pattern:          request.Pattern,
		ExpectedInterval: request.ExpectedInterval,
		Output:           request.Output,
	}
}

func (t *modelTransformer) TransformUpdateRequestToModel(request UpdateRequest) *Heartbeat {
	return t.TransformCreateRequestToModel(CreateRequest{
		BaseEditRequest: request.BaseEditRequest,
		ID:              request.ID,
	})
}

func (t *modelTransformer) TransformBulkCreateRequestToModels(request BulkCreateRequest) []*Heartbeat {
	heartbeats := make([]*Heartbeat, len(request.Items))
	for i := range heartbeats {
		heartbeats[i] = t.TransformCreateRequestToModel(request.Items[i])
	}

	return heartbeats
}

func (t *modelTransformer) TransformBulkUpdateRequestToModels(request BulkUpdateRequest) []*Heartbeat {
	heartbeats := make([]*Heartbeat, len(request.Items))
	for i := range heartbeats {
		heartbeats[i] = t.TransformCreateRequestToModel(CreateRequest{
			BaseEditRequest: request.Items[i].BaseEditRequest,
			ID:              request.Items[i].ID,
		})
	}

	return heartbeats
}
