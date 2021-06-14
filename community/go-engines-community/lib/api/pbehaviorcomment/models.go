package pbehaviorcomment

type Request struct {
	Pbehavior string `json:"pbehavior" binding:"required"`
	Author    string `json:"author" swaggerignore:"true"`
	Message   string `json:"message" binding:"required,max=255"`
}
