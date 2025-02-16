package models

import (
	"errors"
	"testing"

	"github.com/sullyvannunes/todo-list/pkg/models"
)

type Validation struct {
	Field string
	Msg   string
}

func AssertValidation(t *testing.T, v Validation, err error, failMessage string) {
	t.Helper()
	var vErr *models.ValidationErrors
	ok := errors.As(err, &vErr)
	if !ok {
		t.Fatal("must return a validation error")
	}

	var uErr *models.ValidationError
	for _, e := range vErr.Errors {
		if e.Field == v.Field {
			uErr = e
			break
		}
	}

	if uErr == nil {
		t.Fatalf("must validate %s", v.Field)
	}

	for _, msg := range uErr.Msgs {
		if msg == v.Msg {
			return
		}
	}

	t.Fatal(failMessage)
}
