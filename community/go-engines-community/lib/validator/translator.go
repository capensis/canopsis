package validator

import (
	"github.com/go-playground/locales/de"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func NewTranslator(v *validator.Validate) (*ut.UniversalTranslator, error) {
	enTrans := en.New()
	deprecatedTrans := de.New() // todo remove when everything is refactored
	uniTrans := ut.New(enTrans, enTrans, fr.New(), deprecatedTrans)
	err := RegisterTranslations(v, uniTrans, InvalidIDChars)
	depTranslator, _ := uniTrans.GetTranslator("de")
	registerDeprecatedTranslations(v, depTranslator)

	return uniTrans, err
}
