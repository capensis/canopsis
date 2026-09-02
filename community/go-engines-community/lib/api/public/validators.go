package api

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exdate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterValidators() (*validator.Validate, *ut.UniversalTranslator, error) {
	v := validator.New()
	v.SetTagName("binding")
	trans, err := libvalidator.NewTranslator(v, "")
	if err != nil {
		return v, nil, err
	}

	err = registerTagValidations(v)
	if err != nil {
		return v, nil, err
	}

	registerStructValidations(v)

	return v, trans, nil
}

func registerTagValidations(v *validator.Validate) error {
	var tagValidations = map[string]validator.Func{
		"notblank":       validators.NotBlank,
		"oneoforempty":   validation.ValidateOneOfOrEmpty,
		"iscolororempty": validation.ValidateColorOrEmpty,
	}
	for tag, f := range tagValidations {
		err := v.RegisterValidation(tag, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func registerStructValidations(
	v *validator.Validate,
) {
	v.RegisterCustomTypeFunc(validation.ValidateCpsTimeType, datetime.CpsTime{})

	v.RegisterStructValidation(exdate.ValidateExdateRequest, exdate.Request{})

	v.RegisterStructValidation(pbehavior.ValidateEditRequest, pbehavior.EditRequest{})
	v.RegisterStructValidation(pbehavior.ValidatePatchRequest, pbehavior.PatchRequest{})
	v.RegisterStructValidation(pbehavior.ValidateConnectorCreateRequest, pbehavior.BulkConnectorCreateRequestItem{})
	v.RegisterStructValidation(pbehavior.ValidateConnectorEditRequest, pbehavior.BulkConnectorEditRequestItem{})
}
