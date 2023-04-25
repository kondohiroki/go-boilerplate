package validation

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var uni *ut.UniversalTranslator
var trans ut.Translator
var m sync.Mutex

func InitValidator() {
	if validate == nil {
		m.Lock()
		defer m.Unlock()

		// Create a new validator instance
		validate = validator.New()

		// Turns to use json tag instead of struct field name
		// When using validationError.Field() it will return the json tag instead of struct field name
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})

		// Set translator
		en := en.New()
		uni = ut.New(en, en)

		// This is usually know or extracted from http 'Accept-Language' header
		// also see uni.FindTranslator(...)
		trans, _ = uni.GetTranslator("en")

		// Register translation for validator
		en_translations.RegisterDefaultTranslations(validate, trans)

	}
}

func GetValidator() (*validator.Validate, ut.Translator) {
	if validate == nil {
		InitValidator()
	}
	return validate, trans
}

func GetValidationErrors(validationErrors validator.ValidationErrors) (errors []map[string]any) {
	for _, validationError := range validationErrors {
		errorItem := make(map[string]any, 1)
		errorItem[validationError.Field()] = validationError.Translate(trans)
		errors = append(errors, errorItem)
	}

	return errors
}

func Translate(v validator.FieldError) string {
	return v.Translate(trans)
}
