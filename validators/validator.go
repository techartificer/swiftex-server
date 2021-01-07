package validators

import (
	"encoding/json"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	v     *validator.Validate
	trans ut.Translator
)

type ValidationError map[string]string

func (ve *ValidationError) Error() string {
	b, _ := json.Marshal(ve)
	return string(b)
}
func (ve ValidationError) Add(key, value string) {
	ve[key] = value
}

// InitValidator initialize the validator
func InitValidator() {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	v = validator.New()
	en_translations.RegisterDefaultTranslations(v, trans)
}

// GetValidationError returns validation error
func GetValidationError(any interface{}) error {
	ve := ValidationError{}
	if err := v.Struct(any); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			ve.Add(e.Field(), e.Translate(trans))
		}
		return &ve
	}
	return nil
}
