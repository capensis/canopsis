package security

const RoleAdmin = "admin"

const (
	// QueryParamApiKey is the user api key for auth.
	QueryParamApiKey = "authkey"
	// HeaderApiKey is the user api key for auth.
	HeaderApiKey = "x-canopsis-authkey" //nolint:gosec
	// QueryParamCasTicket is CAS ticket for auth.
	QueryParamCasTicket = "ticket"
	// QueryParamCasService is CAS service for auth.
	QueryParamCasService = "service"
	// SessionKey is the session name in cookies.
	SessionKey = "session-id"
)
