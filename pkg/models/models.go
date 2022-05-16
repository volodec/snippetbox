package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: Подходящей записи не найдено")

type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}
