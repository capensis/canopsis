package pbehaviortype

type ModelTransformer interface {
	TransformCreateRequestToModel(request CreateRequest) *Type
	TransformUpdateRequestToModel(request UpdateRequest) *Type
}

func NewModelTransformer() ModelTransformer {
	return &modelTransformer{}
}

type modelTransformer struct{}

func (modelTransformer) TransformCreateRequestToModel(request CreateRequest) *Type {
	return &Type{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Priority:    *request.Priority,
		IconName:    request.IconName,
		Color:       request.Color,
	}
}

func (t *modelTransformer) TransformUpdateRequestToModel(request UpdateRequest) *Type {
	return t.TransformCreateRequestToModel(CreateRequest(request))
}
