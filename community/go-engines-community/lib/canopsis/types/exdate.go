package types

type Exdate struct {
	Begin CpsTime `bson:"begin" json:"begin"`
	End   CpsTime `bson:"end" json:"end"`
}
