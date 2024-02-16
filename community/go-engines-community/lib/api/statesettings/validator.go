package statesettings

import (
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStateSettingRequest(sl validator.StructLevel)
	ValidateStateThresholds(sl validator.StructLevel)
}

type baseValidator struct{}

func NewValidator() Validator {
	return &baseValidator{}
}

func (v *baseValidator) ValidateStateSettingRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(StateSettingRequest)

	switch r.Type {
	case statesetting.TypeJUnit:
		switch r.Method {
		case statesetting.MethodWorst:
			if r.JunitThresholds != nil {
				sl.ReportError(r.JunitThresholds, "junit_thresholds", "junit_thresholds", "must_be_empty", "")
			}
		case statesetting.MethodWorstOfShare:
			if r.JunitThresholds == nil {
				sl.ReportError(r.JunitThresholds, "junit_thresholds", "junit_thresholds", "notblank", "")
			}
		default:
			sl.ReportError(r.Method, "Method", "Method", "oneof", strings.Join([]string{
				statesetting.MethodWorst,
				statesetting.MethodWorstOfShare,
			}, ","))
		}
	default:
		sl.ReportError(r.Method, "Type", "Type", "oneof", strings.Join([]string{
			statesetting.TypeJUnit,
		}, ","))
	}
}

func (v *baseValidator) ValidateStateThresholds(sl validator.StructLevel) {
	r := sl.Current().Interface().(StateThresholds)

	if r.Type != nil {
		switch *r.Type {
		case statesetting.TypeNumber, statesetting.TypePercentage:
			/*do nothing*/
		default:
			sl.ReportError(r.Type, "Type", "Type", "oneof", strings.Join([]string{
				strconv.Itoa(statesetting.TypeNumber),
				strconv.Itoa(statesetting.TypePercentage),
			}, ","))
		}
	}
}
