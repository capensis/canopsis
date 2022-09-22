package types

type Exdate struct {
	Begin CpsTime `bson:"begin" json:"begin" swaggertype:"integer"`
	End   CpsTime `bson:"end" json:"end" swaggertype:"integer"`
}
