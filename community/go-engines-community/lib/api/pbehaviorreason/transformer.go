package pbehaviorreason

type ModelTransformer interface {
	TransformCreateRequestToModel(request CreateRequest) *Reason
	TransformUpdateRequestToModel(request UpdateRequest) *Reason
}

func NewModelTransformer() ModelTransformer {
	return &modelTransformer{}
}

type modelTransformer struct{}

func (modelTransformer) TransformCreateRequestToModel(request CreateRequest) *Reason {
	return &Reason{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
	}
}

func (t *modelTransformer) TransformUpdateRequestToModel(request UpdateRequest) *Reason {
	return t.TransformCreateRequestToModel(CreateRequest(request))
}
