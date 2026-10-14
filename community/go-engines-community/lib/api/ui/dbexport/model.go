package dbexport

type Request struct {
	IDs []string `json:"ids" binding:"required"`
}
