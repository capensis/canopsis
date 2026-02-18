package api

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterValidators(secConfig libsecurity.Config) (*validator.Validate, *ut.UniversalTranslator, error) {
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

	registerStructValidations(v, secConfig)

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

// todo
func registerStructValidations(
	v *validator.Validate,
	secConfig libsecurity.Config,
) {
}
