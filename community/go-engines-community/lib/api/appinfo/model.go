package appinfo

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/colortheme"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type VersionConf struct {
	Edition string `json:"edition" bson:"edition"`
	Stack   string `json:"stack" bson:"stack"`

	Version        string            `json:"version" bson:"version"`
	VersionUpdated *datetime.CpsTime `json:"version_updated" bson:"version_updated" swaggertype:"integer"`
}

type PopupTimeout struct {
	Info  *datetime.DurationWithUnit `json:"info,omitempty" bson:"info,omitempty"`
	Error *datetime.DurationWithUnit `json:"error,omitempty" bson:"error,omitempty"`
}

type UserInterfaceConf struct {
	AppTitle                   string        `json:"app_title,omitempty" bson:"app_title,omitempty"`
	Footer                     string        `json:"footer,omitempty" bson:"footer,omitempty"`
	LoginPageDescription       string        `json:"login_page_description,omitempty" bson:"login_page_description,omitempty"`
	Logo                       string        `json:"logo,omitempty" bson:"logo,omitempty"`
	Language                   string        `json:"language,omitempty" bson:"language,omitempty" binding:"oneoforempty=fr en"`
	PopupTimeout               *PopupTimeout `json:"popup_timeout,omitempty" bson:"popup_timeout,omitempty"`
	AllowChangeSeverityToInfo  bool          `json:"allow_change_severity_to_info" bson:"allow_change_severity_to_info"`
	MaxMatchedItems            int64         `json:"max_matched_items" bson:"max_matched_items" binding:"gt=0"`
	CheckCountRequestTimeout   int64         `json:"check_count_request_timeout" bson:"check_count_request_timeout" binding:"gt=0"`
	ShowHeaderOnKioskMode      bool          `json:"show_header_on_kiosk_mode" bson:"show_header_on_kiosk_mode"`
	RequiredInstructionApprove bool          `json:"required_instruction_approve" bson:"required_instruction_approve"`
}

type GlobalConf struct {
	Timezone          string `json:"timezone,omitempty"`
	FileUploadMaxSize int64  `json:"file_upload_max_size"`

	EventsCountTriggerDefaultThreshold int `json:"events_count_trigger_default_threshold"`
}

type RemediationConf struct {
	JobConfigTypes []JobConfigType `json:"job_config_types"`
}

type JobConfigType struct {
	Name      string `json:"name"`
	AuthType  string `json:"auth_type"`
	WithBody  bool   `json:"with_body"`
	WithQuery bool   `json:"with_query"`
}

type AppInfoResponse struct {
	UserInterfaceConf
	GlobalConf
	VersionConf
	Login       LoginConf        `json:"login"`
	Remediation *RemediationConf `json:"remediation,omitempty"`

	DefaultColorTheme colortheme.Theme `json:"default_color_theme"`
	Maintenance       bool             `json:"maintenance"`

	SerialName string `json:"serial_name"`
}

type LoginConf struct {
	BasicConfig  LoginConfigMethod       `json:"basic,omitempty"`
	CasConfig    LoginConfigMethod       `json:"casconfig,omitempty"`
	LdapConfig   LoginConfigMethod       `json:"ldapconfig,omitempty"`
	SamlConfig   LoginConfigMethod       `json:"saml2config,omitempty"`
	OAuth2Config Oauth2LoginConfigMethod `json:"oauth2config,omitempty"`
}

type LoginConfigMethod struct {
	Title  string `json:"title,omitempty"`
	Enable bool   `json:"enable"`
}

type Oauth2LoginConfigMethod struct {
	Providers []string `json:"providers,omitempty"`
	Enable    bool     `json:"enable"`
}
