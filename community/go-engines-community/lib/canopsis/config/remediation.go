package config

const (
	AuthTypeHeaderToken = "header-token"
	AuthTypeBearerToken = "bearer-token"
	AuthTypeBasicAuth   = "basic-auth"
)

type RemediationConf struct {
	HttpTimeout                    string                       `toml:"http_timeout" bson:"http_timeout"`
	LaunchJobRetriesAmount         int                          `toml:"launch_job_retries_amount" bson:"launch_job_retries_amount"`
	LaunchJobRetriesInterval       string                       `toml:"launch_job_retries_interval" bson:"launch_job_retries_interval"`
	WaitJobCompleteRetriesAmount   int                          `toml:"wait_job_complete_retries_amount" bson:"wait_job_complete_retries_amount"`
	WaitJobCompleteRetriesInterval string                       `toml:"wait_job_complete_retries_interval" bson:"wait_job_complete_retries_interval"`
	PauseManualInstructionInterval string                       `toml:"pause_manual_instruction_interval" bson:"pause_manual_instruction_interval"`
	ExternalAPI                    map[string]ExternalApiConfig `toml:"external_api" bson:"external_api"`
}

// ExternalApiConfig represents configuration of external service API.
type ExternalApiConfig struct {
	Auth              Auth           `toml:"auth" bson:"auth"`
	LaunchEndpoint    LaunchEndpoint `toml:"launch_endpoint" bson:"launch_endpoint"`
	StatusEndpoint    StatusEndpoint `toml:"status_endpoint" bson:"status_endpoint"`
	QueueEndpoint     QueueEndpoint  `toml:"queue_endpoint" bson:"queue_endpoint"`
	OutputEndpoint    OutputEndpoint `toml:"output_endpoint" bson:"output_endpoint"`
	ResponseErrMsgKey string         `toml:"response_err_msg_key" bson:"response_err_msg_key,omitempty"`
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
