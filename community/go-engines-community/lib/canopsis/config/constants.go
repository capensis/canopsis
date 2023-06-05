package config

import (
	"html/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// Default config values
const (
	AlarmCancelAutosolveDelay = 60 * 60 * time.Second
	AlarmDefaultNameScheme    = "{{ rand_string 2 }}-{{ rand_string 2 }}-{{ rand_string 2 }}"

	ApiTokenSigningMethod = "HS256"
	ApiBulkMaxSize        = 10000

	RemediationHttpTimeout                    = 30 * time.Second
	RemediationLaunchJobRetriesAmount         = 3
	RemediationLaunchJobRetriesInterval       = 5 * time.Second
	RemediationWaitJobCompleteRetriesAmount   = 12
	RemediationWaitJobCompleteRetriesInterval = 5 * time.Second
	RemediationPauseManualInstructionInterval = 15 * time.Second

	MetricsFlushInterval  = 10 * time.Second
	MetricsSliInterval    = time.Hour
	MetricsMaxSliInterval = time.Hour

	TechMetricsDumpKeepInterval = time.Hour

	UserInterfaceMaxMatchedItems          = 10000
	UserInterfaceCheckCountRequestTimeout = 30

	DataStorageMaxUpdates = 100000
)

func CreateDisplayNameTpl(text string) (*template.Template, error) {
	return template.New("displayname_gen_scheme").
		Funcs(template.FuncMap{
			"rand_string": func(n int) string {
				return utils.RandString(n)
			},
		}).
		Parse(text)
}
