package appinfo

const (
	Casconfig   = "casconfig"
	Ldapconfig  = "ldapconfig"
	CrecordName = "crecord_name"
)

type LoginServiceConf struct {
	CrecordName string                 `json:"crecord_name" bson:"crecord_name"`
	Enable      bool                   `json:"enable" bson:"enable"`
	Fields      map[string]interface{} `bson:",inline"`
}

type CanopsisVersionConf struct {
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
	MaxMatchedItems           int64         `json:"max_matched_items" bson:"max_matched_items"`
	CheckCountRequestTimeout  int64         `json:"check_count_request_timeout" bson:"check_count_request_timeout"`
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
	CanopsisVersionConf
	RemediationConf
}

type LogInProvider struct {
}

type LoginConfig struct {
	Providers  map[string]int `json:"providers,omitempty"`
	Casconfig  interface{}    `json:"casconfig,omitempty"`
	Ldapconfig *struct {
		Enable bool `json:"enable"`
	} `json:"ldapconfig,omitempty"`
}

type LoginConfigResponse struct {
	LoginConfig       LoginConfig       `json:"login_config"`
	UserInterfaceConf UserInterfaceConf `json:"user_interface"`
	CanopsisVersionConf
}