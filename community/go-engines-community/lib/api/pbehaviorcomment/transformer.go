package pbehaviorcomment

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
)

type ModelTransformer interface {
	TransformRequestToModel(request *Request) *pbehavior.Comment
}

type modelTransformer struct{}

func NewModelTransformer() ModelTransformer {
	return &modelTransformer{}
}

func (t *modelTransformer) TransformRequestToModel(request *Request) *pbehavior.Comment {
	return &pbehavior.Comment{
		Author:  request.Author,
		Message: request.Message,
	}
}
