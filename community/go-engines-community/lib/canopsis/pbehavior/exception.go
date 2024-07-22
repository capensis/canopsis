package pbehavior

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Exdate struct {
	types.Exdate `bson:"inline"`
	Type         string `bson:"type" json:"type"`
}

type Exception struct {
	ID          string           `bson:"_id,omitempty" json:"_id"`
	Name        string           `bson:"name" json:"name"`
	Author      string           `bson:"author" json:"author"`
	Description string           `bson:"description" json:"description"`
	Exdates     []Exdate         `bson:"exdates" json:"exdates"`
	Created     datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty"`
	Updated     datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty"`
}
