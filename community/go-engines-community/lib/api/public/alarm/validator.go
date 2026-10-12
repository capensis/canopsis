package alarm

import (
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
)

var alarmStepTypes = types.GetAlarmStepTypes()

func ValidateStepRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(StepRequest)
	errParam := strings.Join(alarmStepTypes, " ")
	for i, t := range r.Types {
		if !slices.Contains(alarmStepTypes, t) {
			sl.ReportError(t, "Types."+strconv.Itoa(i), "Types."+strconv.Itoa(i), "oneof", errParam)
		}
	}
}
