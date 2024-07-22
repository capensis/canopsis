package cas

type casLoginRequest struct {
	// Redirect is front-end url to redirect back after authentication.
	Redirect string `form:"redirect"`
	// Service is proxy url to callback handler to set auth for front-end app.
	Service string `form:"service"`
}
