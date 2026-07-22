package amqp

// Environment variables linked to session parameters
const (
	EnvURL                    = "CPS_AMQP_URL"
	EnvHTTPURL                = "CPS_AMQP_HTTP_URL"
	EnvHTTPUser               = "CPS_AMQP_HTTP_USER"
	EnvHTTPPassword           = "CPS_AMQP_HTTP_PASSWORD" //nolint:gosec
	EnvHTTPInsecureSkipVerify = "CPS_AMQP_HTTP_INSECURE_SKIP_VERIFY"
)
