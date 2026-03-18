package model

import (
	"app/internal/crosscutting/timer"
	"time"
)

type HasRequestAt struct {
	RequestAt *time.Time
}

func (model HasRequestAt) GetRequestAt() time.Time {
	if model.RequestAt == nil {
		return timer.Now()
	}
	return *model.RequestAt
}

func (model *HasRequestAt) SetRequestAt(requestAt time.Time) {
	model.RequestAt = &requestAt
}
