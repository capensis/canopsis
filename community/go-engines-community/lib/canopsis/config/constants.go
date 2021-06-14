package config

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"html/template"
	"time"
)

// Default config values
const (
	AlarmBaggotTime                    = 60 * time.Second
	AlarmCancelAutosolveDelay          = 60 * 60 * time.Second
	AlarmDefaultNameScheme             = "{{ rand_string 2 }}-{{ rand_string 2 }}-{{ rand_string 2 }}"
	RemediationJobExecutorFetchTimeout = 30 * time.Second
	UserInterfaceMaxPbhEntities        = 10000
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
