package auth

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"` //nolint:gosec
}

type LoginResponse struct {
	AccessToken string `json:"access_token"` //nolint:gosec
}

type LoggedUserCountResponse struct {
	Count int64 `json:"count"`
}
