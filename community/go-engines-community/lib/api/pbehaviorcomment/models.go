package pbehaviorcomment

type Request struct {
	Pbehavior string `json:"pbehavior" binding:"required"`
	Author    string `json:"author" binding:"required,max=255"`
	Message   string `json:"message" binding:"required,max=255"`
}
