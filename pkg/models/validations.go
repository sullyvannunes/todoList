package models

import (
	"context"
	"reflect"
)

type ValidationFunc func(context.Context) *ValidationError

func ValidatePresenceOf(field string, obj any) ValidationFunc {
	return func(ctx context.Context) *ValidationError {
		v := reflect.ValueOf(obj)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return nil
		}

		f := v.FieldByName(field)
		if f.IsZero() {
			return &ValidationError{Field: field, Msgs: []string{"Must not be blank"}}
		}

		return nil
	}
}

func Validate(ctx context.Context, funcs ...ValidationFunc) error {
	errs := &ValidationErrors{}

	for _, f := range funcs {
		err := f(ctx)
		if err != nil {
			errs.Errors = append(errs.Errors, err)
		}
	}

	if len(errs.Errors) == 0 {
		return nil
	}

	return errs
}
