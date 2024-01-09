package file

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type File struct {
	ID        string           `bson:"_id" json:"_id"`
	FileName  string           `bson:"filename" json:"filename"`
	MediaType string           `bson:"mediatype" json:"mediatype"`
	Created   datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Storage   string           `bson:"storage" json:"-"`
	Etag      string           `bson:"etag" json:"-"`
	IsPublic  bool             `bson:"is_public" json:"-"`
}

type CreateRequest struct {
	Public bool `form:"public" json:"public"`
}
