package bulk

type DeleteRequestItem struct {
	ID     string `json:"_id" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
}
