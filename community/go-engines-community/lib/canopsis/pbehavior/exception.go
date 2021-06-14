package pbehavior

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
)

const (
	ExceptionCollectionName = mongo.PbehaviorExceptionMongoCollection
)

type Exdate struct {
	Begin types.CpsTime `bson:"begin" json:"begin"`
	End   types.CpsTime `bson:"end" json:"end"`
	Type  string        `bson:"type" json:"type"`
}

type Exception struct {
	ID          string         `bson:"_id,omitempty" json:"_id"`
	Name        string         `bson:"name" json:"name"`
	Description string         `bson:"description" json:"description"`
	Exdates     []Exdate       `bson:"exdates" json:"exdates"`
	Created     *types.CpsTime `bson:"created,omitempty" json:"created"`
}
