package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Api       ApiConfig `yaml:"api"`
	Pbehavior struct {
		Type   string `yaml:"type"`
		Reason string `yaml:"reason"`
	} `yaml:"pbehavior"`
	Timezone string           `yaml:"timezone"`
	Inactive []InactiveConfig `yaml:"inactive"`
	location *time.Location
}

type ApiConfig struct {
	Host               string `yaml:"host"`
	InsecureSkipverify bool   `yaml:"insecure_skip_verify"`

	Username string `yaml:"-"`
	Password string `yaml:"-"`
}

type HourRange struct {
	hhmm     [2]int
	weekdays [7]bool
}

type InactiveConfig struct {
	Hours      []string
	Hostgroups []string
	hourRanges []HourRange
}

// LoadConfig reads a file in configPath path and parses its YAML content.
func LoadConfig(configPath string) (Config, error) {
	config := Config{}
	buf, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return config, err
	}

	if config.Timezone != "" {
		config.location, err = time.LoadLocation(config.Timezone)
		if err != nil {
			return config, err
		}
	}
	if config.location == nil {
		config.location, err = time.LoadLocation(defaultLocation)
		if err != nil {
			return config, err
		}
	}
	for i, c := range config.Inactive {
		if err = c.ParseHours(); err != nil {
			return config, err
		}
		config.Inactive[i] = c
	}

	return config, err
}

func hhmm2i(s string) (int, error) {
	sh, sm, _ := strings.Cut(s, ":")
	h, err := strconv.Atoi(sh)
	if err != nil {
		return 0, err
	}
	m := 0
	if sm != "" {
		m, err = strconv.Atoi(sm)
		if err != nil {
			return 0, err
		}
	}
	return h*100 + m, nil
}

func idx(s string, sl []string) int {
	for i := 0; i < len(sl); i++ {
		if s == sl[i] {
			return i
		}
	}
	return -1
}

func (c ApiConfig) CreateRequest(ctx context.Context, method, path string, b []byte, q url.Values) (*http.Request, error) {
	u, err := url.Parse(c.Host)
	if err != nil {
		return nil, fmt.Errorf("host %q is invalid url: %w", c.Host, err)
	}
	u.Path = path
	if q != nil {
		u.RawQuery = q.Encode()
	}

	var body io.Reader
	if len(b) > 0 {
		body = bytes.NewReader(b)
	}
	request, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", canopsis.JsonContentType)
	request.Header.Set("Accept", canopsis.JsonContentType)

	request.SetBasicAuth(c.Username, c.Password)

	return request, nil
}

// ParseHours translates Hours rules. Each rule can be as "19:30-5;TU,TH-MO" with local start-end times range
// and weekdays separated by comma as list or range(s)
func (inactive *InactiveConfig) ParseHours() error {
	const weekLength = 7
	weekdayNames := [weekLength]string{"MO", "TU", "WE", "TH", "FR", "SA", "SU"}

	for _, h := range inactive.Hours {
		hours, sWeekdays, _ := strings.Cut(h, ";")
		sStart, sEnd, ok := strings.Cut(hours, "-")
		if !ok {
			return fmt.Errorf("invalid value of time range %q", h)
		}
		start, err := hhmm2i(sStart)
		if err != nil {
			return err
		}
		end, err := hhmm2i(sEnd)
		if err != nil {
			return err
		}
		weekdays := [7]bool{}
		if sWeekdays != "" {
			weekIntervals := strings.Split(sWeekdays, ",")
			for _, interval := range weekIntervals {
				sStart, sEnd, ok := strings.Cut(interval, "-")
				idxStart := idx(sStart, weekdayNames[:])
				if idxStart == -1 {
					continue
				}
				if !ok {
					weekdays[idxStart] = true
					continue
				}
				idxEnd := idx(sEnd, weekdayNames[:])
				if idxEnd < idxStart {
					idxEnd += weekLength
				}
				for i := idxStart; i <= idxEnd; i++ {
					if i+1 > weekLength {
						weekdays[i-weekLength] = true
					} else {
						weekdays[i] = true
					}
				}
			}
		} else {
			for i := 0; i < 7; i++ {
				weekdays[i] = true
			}
		}

		inactive.hourRanges = append(inactive.hourRanges, HourRange{
			hhmm:     [2]int{start, end},
			weekdays: weekdays,
		})
	}
	return nil
}
