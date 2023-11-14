package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EntityService struct {
	types.Entity   `bson:",inline"`
	OutputTemplate string `bson:"output_template" json:"output_template"`

	savedpattern.EntityPatternFields `bson:",inline"`
}
