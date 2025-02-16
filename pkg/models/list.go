package models

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

type ListId ulid.ULID
type List struct {
	Id        ListId
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    UserId
}

func (m *Models) CreateList(ctx context.Context, l *List) error {
	err := Validate(ctx, ValidatePresenceOf("UserId", l), ValidatePresenceOf("Name", l))
	if err != nil {
		return err
	}
	return nil
}
