package validator

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func NewTranslator(v *validator.Validate, invalidIDChars string) (*ut.UniversalTranslator, error) {
	enTrans := en.New()
	uniTrans := ut.New(enTrans, enTrans, fr.New())
	err := RegisterTranslations(v, uniTrans, invalidIDChars)

	return uniTrans, err
}
