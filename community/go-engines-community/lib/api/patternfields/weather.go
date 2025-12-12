package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type WeatherServiceRequest struct {
	WeatherServicePattern          pattern.WeatherServicePattern `json:"weather_service_pattern" binding:"weather_service_pattern"`
	CorporateWeatherServicePattern string                        `json:"corporate_weather_service_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
}

func (r WeatherServiceRequest) ToModel() savedpattern.WeatherServicePatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.WeatherServicePatternFields{
			WeatherServicePattern: r.WeatherServicePattern,
		}
	}

	return savedpattern.WeatherServicePatternFields{
		WeatherServicePattern:               r.CorporatePattern.WeatherServicePattern,
		CorporateWeatherServicePattern:      r.CorporatePattern.ID,
		CorporateWeatherServicePatternTitle: r.CorporatePattern.Title,
	}
}
