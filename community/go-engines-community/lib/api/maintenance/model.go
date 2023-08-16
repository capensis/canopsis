package maintenance

type Request struct {
	Message string `bson:"message"`
	Enabled *bool  `bson:"enabled" binding:"required"`
}
