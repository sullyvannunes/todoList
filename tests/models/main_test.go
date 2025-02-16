package models

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/sullyvannunes/todo-list/pkg/models"
)

var r *MockRepository
var mo *models.Models

func TestMain(m *testing.M) {
	r = &MockRepository{
		db: make(map[ulid.ULID]interface{}),
	}
	mo = models.NewModels(r)
	m.Run()
}
