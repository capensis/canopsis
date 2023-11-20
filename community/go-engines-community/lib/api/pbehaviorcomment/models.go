package pbehaviorcomment

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
)

type Request struct {
	Pbehavior string `json:"pbehavior" binding:"required"`
	Author    string `json:"author" swaggerignore:"true"`
	Message   string `json:"message" binding:"required,max=255"`
}

type Response struct {
	ID        string          `bson:"_id" json:"_id"`
	Author    *author.Author  `bson:"author" json:"author"`
	Timestamp libtime.CpsTime `bson:"ts" json:"ts" swaggertype:"integer"`
	Message   string          `bson:"message" json:"message"`
}
