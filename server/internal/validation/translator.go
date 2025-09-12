package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translation "github.com/go-playground/validator/v10/translations/en"
)

var Transl ut.Translator

func SetupTranslations(customValidator *CustomValidator) {
	en := en.New()
	unt := ut.New(en, en)
	Transl, _ = unt.GetTranslator("en")
	en_translation.RegisterDefaultTranslations(customValidator.validator, Transl)
}
