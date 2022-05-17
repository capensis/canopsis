package pbehavior

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

const (
	ReasonCollectionName = mongo.PbehaviorReasonMongoCollection
)

type Reason struct {
	ID          string        `bson:"_id,omitempty" json:"_id"`
	Name        string        `bson:"name" json:"name" binding:"required"`
	Description string        `bson:"description" json:"description" binding:"required"`
	Created     types.CpsTime `bson:"created,omitempty" json:"created" swaggertype:"integer"`
}
