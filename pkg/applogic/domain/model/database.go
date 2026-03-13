package model

import (
	"time"
)

type HasCreatedAt struct {
	CreatedAt time.Time
}

type HasUpdatedAt struct {
	UpdatedAt time.Time
}

type DBItem struct {
	HasCreatedAt
	HasUpdatedAt
}

// to implement DBRecorder interface
func (model DBItem) DBRecordable() {}

func (model *DBItem) Create(t time.Time) {
	model.CreatedAt = t
	model.UpdatedAt = t
}

func (model *DBItem) Update(t time.Time) {
	model.UpdatedAt = t
}

type DBRecorder interface {
	DBRecordable()
}

type HasName struct {
	Name string
}
