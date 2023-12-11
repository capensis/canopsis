package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

type EntityService struct {
	types.Entity   `bson:",inline"`
	OutputTemplate string `bson:"output_template" json:"output_template"`

	savedpattern.EntityPatternFields `bson:",inline"`
}

func (s *EntityService) GetMongoQueries() (bson.M, bson.M, error) {
	var query, negativeQuery bson.M
	var err error

	if len(s.EntityPattern) > 0 {
		query, err = db.EntityPatternToMongoQuery(s.EntityPattern, "")
		if err != nil {
			return nil, nil, err
		}

		negativeQuery, err = db.EntityPatternToNegativeMongoQuery(s.EntityPattern, "")
		if err != nil {
			return nil, nil, err
		}
	}

	return query, negativeQuery, nil
}
