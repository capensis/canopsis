package config

import (
	"html/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// Default config values
const (
	AlarmCancelAutosolveDelay = 60 * 60 * time.Second
	AlarmDisplayNameScheme    = "{{ rand_string 2 }}-{{ rand_string 2 }}-{{ rand_string 2 }}"

	ApiTokenSigningMethod               = "HS256"
	ApiBulkMaxSize                      = 10000
	ApiExportMongoClientTimeout         = time.Minute
	ApiMetricsCacheExpiration           = 24 * time.Hour
	ApiEventsRecorderFetchStatusTimeout = 3 * time.Second

	RemediationHttpTimeout                    = 5 * time.Second
	RemediationPauseManualInstructionInterval = 15 * time.Second
	RemediationJobRetryInterval               = 10 * time.Second
	RemediationJobWaitInterval                = 60 * time.Second

	MetricsFlushInterval          = 10 * time.Second
	MetricsSliInterval            = time.Hour
	MetricsMaxSliInterval         = time.Hour
	MetricsUserSessionGapInterval = time.Hour

	TechMetricsDumpKeepInterval  = time.Hour
	TechMetricsGoMetricsInterval = time.Second

	UserInterfaceMaxMatchedItems          = 10000
	UserInterfaceCheckCountRequestTimeout = 30

	DataStorageMaxUpdates = 100000

	DefaultEventsCountThreshold = 10
)

var ApiAuthorScheme = []string{"$username"}

func CreateDisplayNameTpl(text string) (*template.Template, error) {
	return template.New("displayname_gen_scheme").
		Funcs(template.FuncMap{
			"rand_string": func(n int) string {
				return utils.RandString(n)
			},
		}).
		Parse(text)
}
