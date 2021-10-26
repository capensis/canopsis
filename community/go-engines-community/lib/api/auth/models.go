package auth

type loginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

type loggedUserCountResponse struct {
	Count int64 `json:"count"`
}
