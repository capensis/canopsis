package pattern

type WeatherServicePattern [][]FieldCondition

func (p WeatherServicePattern) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}
