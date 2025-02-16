package models

import (
	"context"
	"testing"

	"github.com/sullyvannunes/todo-list/pkg/models"
)

func TestList(t *testing.T) {
	ctx := context.Background()

	t.Run("Validates user presence", func(t *testing.T) {
		l := &models.List{}
		err := mo.CreateList(ctx, l)
		AssertValidation(t, Validation{"UserId", "Must not be blank"}, err, "must validate user presence")
	})

	t.Run("Validates list name", func(t *testing.T) {
		l := &models.List{}
		err := mo.CreateList(ctx, l)
		AssertValidation(t, Validation{"Name", "Must not be blank"}, err, "must validate name presence")
	})

	// t.Run("Create a list", func(t *testing.T) {
	// 	l := &models.List{Name: "First List", UserId: models.UserId(ulid.Make())}
	// 	err := mo.CreateList(ctx, l)
	// 	assert.NoError(t, err)
	// })
}
