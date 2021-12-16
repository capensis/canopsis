package widgetfilter

type CreateRequest struct {
	EditRequest
	Widget   string `json:"widget" binding:"required"`
	Personal *bool  `json:"personal" binding:"required"`
}

type EditRequest struct {
	ID     string `json:"-"`
	Title  string `json:"title" binding:"required,max=255"`
	Query  string `json:"query" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
}
