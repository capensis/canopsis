package userpreferences

type EditRequest struct {
	Content map[string]interface{} `json:"content" binding:"required"`
}

type Response struct {
	Widget  string                 `bson:"widget" json:"widget"`
	Content map[string]interface{} `bson:"content" json:"content"`
}
