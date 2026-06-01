package validation

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Struct  string `json:"struct"`
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func Format(err error) []FieldError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return []FieldError{}
	}

	out := make([]FieldError, 0, len(ve))
	for _, fe := range ve {
		out = append(out, FieldError{
			Struct:  fe.StructNamespace(),
			Field:   fe.Field(),
			Rule:    fe.Tag(),
			Message: humanize(fe),
		})
	}

	return out
}

func AsValidationErrors(err error, out *validator.ValidationErrors) bool {
	if err == nil {
		return false
	}
	return errors.As(err, out)
}

func humanize(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
		}
		return fmt.Sprintf("%s must be >= %s", fe.Field(), fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
		}
		return fmt.Sprintf("%s must be <= %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s failed validation (%s)", fe.Field(), fe.Tag())
	}
}
