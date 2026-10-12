package sessionauth

type loginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"` //nolint:gosec
}

type loginResponse struct {
	Name       string   `json:"crecord_name"`
	AuthApiKey string   `json:"authkey"` //nolint:gosec
	Roles      []string `json:"roles"`
	Contact    struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"contact"`
	Email string `json:"mail"`
}
