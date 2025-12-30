package validator

import (
	"reflect"
	"strings"

	v "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validate *v.Validate
}

func New() *CustomValidator {
	val := v.New()

	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validate: val}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}
