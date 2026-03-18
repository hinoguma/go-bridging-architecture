package infra

import (
	"app/internal/crosscutting"
	"github.com/google/uuid"
)

type UUIDV4Generator struct {
}

func (generator UUIDV4Generator) Issue() string {
	return uuid.New().String()
}

func NewUUIDV4Generator() crosscutting.StrIDGenerator {
	return UUIDV4Generator{}
}
