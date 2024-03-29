package config

const (
	AuthTypeHeaderToken = "header-token"
	AuthTypeBearerToken = "bearer-token"
	AuthTypeBasicAuth   = "basic-auth"
)

type RemediationConf struct {
	HttpTimeout                    string                       `toml:"http_timeout" bson:"http_timeout"`
	PauseManualInstructionInterval string                       `toml:"pause_manual_instruction_interval" bson:"pause_manual_instruction_interval"`
	ExternalAPI                    map[string]ExternalApiConfig `toml:"external_api" bson:"external_api"`
	JobWaitInterval                string                       `toml:"job_wait_interval" bson:"job_wait_interval"`
	JobRetryInterval               string                       `toml:"job_retry_interval" bson:"job_retry_interval"`
}

// ExternalApiConfig represents configuration of external service API.
type ExternalApiConfig struct {
	Auth              Auth           `toml:"auth" bson:"auth"`
	ResetEndpoint     ResetEndpoint  `toml:"reset_endpoint" bson:"reset_endpoint"`
	LaunchEndpoint    LaunchEndpoint `toml:"launch_endpoint" bson:"launch_endpoint"`
	StatusEndpoint    StatusEndpoint `toml:"status_endpoint" bson:"status_endpoint"`
	QueueEndpoint     QueueEndpoint  `toml:"queue_endpoint" bson:"queue_endpoint"`
	OutputEndpoint    OutputEndpoint `toml:"output_endpoint" bson:"output_endpoint"`
	ErrOutputEndpoint OutputEndpoint `toml:"err_output_endpoint" bson:"err_output_endpoint"`
	ResponseErrMsgKey string         `toml:"response_err_msg_key" bson:"response_err_msg_key,omitempty"`
	SkipVerify        bool           `toml:"-" bson:"-"`
}

// ResetEndpoint represents API endpoint to reset external job.
type ResetEndpoint struct {
	UrlTpl  string            `toml:"url_tpl" bson:"url_tpl"`
	Method  string            `toml:"method" bson:"method"`
	Headers map[string]string `toml:"headers" bson:"headers,omitempty"`
	Body    string            `toml:"body" bson:"body"`
}

// LaunchEndpoint represents API endpoint to launch external job.
type LaunchEndpoint struct {
	UrlTpl                     string            `toml:"url_tpl" bson:"url_tpl"`
	UrlTplWithParams           string            `toml:"url_tpl_with_params" bson:"url_tpl_with_params"`
	Method                     string            `toml:"method" bson:"method"`
	Headers                    map[string]string `toml:"headers" bson:"headers,omitempty"`
	ResponseStatusUrlKey       string            `toml:"response_status_url_key" bson:"response_status_url_key,omitempty"`
	ResponseStatusHeaderUrlKey string            `toml:"response_status_header_url_key" bson:"response_status_header_url_key,omitempty"`
	ResponseExternalUrlKey     string            `toml:"response_external_url_key" bson:"response_external_url_key,omitempty"`
	Body                       string            `toml:"body" bson:"body"`
	WithBody                   bool              `toml:"with_body" bson:"with_body"`
	WithUrlQuery               bool              `toml:"with_url_query" bson:"with_url_query"`
}

// QueueEndpoint represents API endpoint to fetch external job execution.
type QueueEndpoint struct {
	UrlTpl                 string            `toml:"url_tpl" bson:"url_tpl"`
	Method                 string            `toml:"method" bson:"method"`
	Headers                map[string]string `toml:"headers" bson:"headers,omitempty"`
	ResponseStatusUrlKey   string            `toml:"response_status_url_key" bson:"response_status_url_key,omitempty"`
	ResponseExternalUrlKey string            `toml:"response_external_url_key" bson:"response_external_url_key,omitempty"`
}

// StatusEndpoint represents API endpoint to fetch execution status of external job.
type StatusEndpoint struct {
	UrlTpl            string            `toml:"url_tpl" bson:"url_tpl"`
	Method            string            `toml:"method" bson:"method"`
	Headers           map[string]string `toml:"headers" bson:"headers,omitempty"`
	ResponseIdKey     string            `toml:"response_id_key" bson:"response_id_key,omitempty"`
	ResponseStatusKey string            `toml:"response_status_key" bson:"response_status_key,omitempty"`
	Statuses          map[string]string `toml:"statuses" bson:"statuses"`
}

// OutputEndpoint represents API endpoint to fetch execution output of external job.
type OutputEndpoint struct {
	UrlTpl  string            `toml:"url_tpl" bson:"url_tpl"`
	Method  string            `toml:"method" bson:"method"`
	Headers map[string]string `toml:"headers" bson:"headers,omitempty"`
}

// Auth represents configuration for API auth of external service API.
type Auth struct {
	Type string `toml:"type" bson:"type"`
	// Header is only AuthTypeHeaderToken
	Header string `toml:"header" bson:"header,omitempty"`
}
