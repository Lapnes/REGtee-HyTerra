package validator

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorDetails map[string]string

func (details *ErrorDetails) Add(key, val string) {
	if (*details)[key] != "" {
		(*details)[key] += " | "
	}
	(*details)[key] += val
}

var Validator = func() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}()

func validateProcessable(data interface{}) (details *ErrorDetails, code int) {
	details = &ErrorDetails{}
	code = http.StatusOK

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldT := t.Field(i)

		name := strings.SplitN(fieldT.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			name = strings.ToLower(fieldT.Name)
		}

		tags := fieldT.Tag.Get("process")
		err := Validator.Var(field.Interface(), tags)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				details.Add(name, err.Tag())
			}
			code = http.StatusUnprocessableEntity
		}
	}

	return details, code
}

func Validate(data interface{}) (details *ErrorDetails, code int) {
	details = &ErrorDetails{}
	code = http.StatusOK

	err := Validator.Struct(data)
	if err != nil {
		if errV, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errV {
				details.Add(err.Field(), err.Tag())
			}
			code = http.StatusBadRequest
		}
	}

	prDet, prCode := validateProcessable(data)
	if code != http.StatusBadRequest {
		if prCode != http.StatusOK {
			code = prCode
			for field, det := range *prDet {
				details.Add(field, det)
			}
		}
	}

	if code == http.StatusOK {
		details = nil
	}

	return details, code
}
