package appinfo

type VersionConf struct {
	Edition string `json:"edition" bson:"edition"`
	Stack   string `json:"stack" bson:"stack"`
	Version string `json:"version" bson:"version"`
}

type IntervalUnit struct {
	Interval uint   `json:"interval" bson:"interval"`
	Unit     string `json:"unit" bson:"unit" binding:"oneof=s h m"`
}

type PopupTimeout struct {
	Info  *IntervalUnit `json:"info,omitempty" bson:"info,omitempty"`
	Error *IntervalUnit `json:"error,omitempty" bson:"error,omitempty"`
}

type UserInterfaceConf struct {
	AppTitle                  string        `json:"app_title,omitempty" bson:"app_title,omitempty"`
	Footer                    string        `json:"footer,omitempty" bson:"footer,omitempty"`
	LoginPageDescription      string        `json:"login_page_description,omitempty" bson:"login_page_description,omitempty"`
	Logo                      string        `json:"logo,omitempty" bson:"logo,omitempty"`
	Language                  string        `json:"language,omitempty" bson:"language,omitempty" binding:"oneoforempty=fr en"`
	PopupTimeout              *PopupTimeout `json:"popup_timeout,omitempty" bson:"popup_timeout,omitempty"`
	AllowChangeSeverityToInfo bool          `json:"allow_change_severity_to_info" bson:"allow_change_severity_to_info"`
	MaxMatchedItems           int64         `json:"max_matched_items" bson:"max_matched_items" binding:"gt=0"`
	CheckCountRequestTimeout  int64         `json:"check_count_request_timeout" bson:"check_count_request_timeout" binding:"gt=0"`
}

type TimezoneConf struct {
	Timezone string `json:"timezone,omitempty" bson:"timezone"`
}

type RemediationConf struct {
	JobExecutorFetchTimeoutSeconds int64 `json:"jobexecutorfetchtimeoutseconds,omitempty" bson:"jobexecutorfetchtimeoutseconds"`
}

type AppInfoResponse struct {
	UserInterfaceConf
	TimezoneConf
	VersionConf
	RemediationConf
}

type LoginConfig struct {
	CasConfig  LoginConfigMethod `json:"casconfig,omitempty"`
	LdapConfig LoginConfigMethod `json:"ldapconfig,omitempty"`
	SamlConfig LoginConfigMethod `json:"saml2config,omitempty"`
}

type LoginConfigMethod struct {
	Title  string `json:"title,omitempty"`
	Enable bool   `json:"enable"`
}

type LoginConfigResponse struct {
	LoginConfig       LoginConfig       `json:"login_config"`
	UserInterfaceConf UserInterfaceConf `json:"user_interface"`
	VersionConf
}
