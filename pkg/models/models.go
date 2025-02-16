package models

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type ValidationError struct {
	Field string
	Msgs  []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s %v", v.Field, v.Msgs)
}

type ValidationErrors struct {
	Errors []*ValidationError
}

func (v *ValidationErrors) Error() string {
	s := []string{}
	for _, e := range v.Errors {
		s = append(s, e.Error())
	}

	return fmt.Sprintf("%v", s)
}

type Repository interface {
	CreateList(context.Context, *List) error
	GetListById(context.Context, ListId) (*List, error)
	CreateUser(context.Context, *User) error
	GetUserById(context.Context, UserId) (*User, error)
}

type Models struct {
	db Repository
}

func NewModels(db Repository) *Models {
	return &Models{
		db: db,
	}
}

type UserId ulid.ULID
type User struct {
	Id             UserId
	Name           string
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastLoggedInAt time.Time
}

type ActionId ulid.ULID
type Action struct {
	Id        ActionId
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
