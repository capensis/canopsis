package template

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
)

func ValidateEditDataRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditDataRequest)
	if r.Type != nil && *r.Type == TypeEvent && r.Body != nil {
		b, err := json.Marshal(r.Body)
		if err == nil {
			event := types.Event{}
			err = json.Unmarshal(b, &event)
			if err == nil {
				err = event.InjectExtraInfos(b)
			}
		}

		if err != nil {
			sl.ReportError(r.Body, "Body", "Body", "invalid", "")
		}
	}
}
