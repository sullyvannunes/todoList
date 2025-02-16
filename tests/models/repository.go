package models

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sullyvannunes/todo-list/pkg/models"
)

type MockRepository struct {
	db map[ulid.ULID]interface{}
}

func (r *MockRepository) CreateList(_ context.Context, l *models.List) error {
	u := ulid.Make()
	l.Id = models.ListId(u)
	l.CreatedAt, l.UpdatedAt = time.Now(), time.Now()
	r.db[u] = l
	return nil
}

func (r *MockRepository) GetListById(_ context.Context, id models.ListId) (*models.List, error) {
	l, exist := r.db[ulid.ULID(id)]
	if !exist {
		return nil, fmt.Errorf("list not found")
	}

	list, ok := l.(*models.List)
	if !ok {
		return nil, fmt.Errorf("not a list")
	}

	return list, nil
}

func (r *MockRepository) CreateUser(_ context.Context, us *models.User) error {
	u := ulid.Make()
	us.Id = models.UserId(u)
	us.CreatedAt, us.UpdatedAt = time.Now(), time.Now()
	r.db[u] = us
	return nil
}

func (r *MockRepository) GetUserById(_ context.Context, id models.UserId) (*models.User, error) {
	u, exist := r.db[ulid.ULID(id)]
	if !exist {
		return nil, fmt.Errorf("user not found")
	}

	user, ok := u.(*models.User)
	if !ok {
		return nil, fmt.Errorf("not a user")
	}

	return user, nil
}
