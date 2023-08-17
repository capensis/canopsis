package maintenance

type Request struct {
	Message string `json:"message"`
	Color   string `json:"color" binding:"iscolororempty"`
	Enabled *bool  `json:"enabled" binding:"required"`
}
