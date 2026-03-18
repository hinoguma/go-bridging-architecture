package errors

import (
	"errors"

	serrors "github.com/hinoguma/go-structured-error"
)

const TypeDataNotFound serrors.ErrorType = "INVALID_PASSWORD"

func IsDataNotFoundError(err error) bool {
	return serrors.IsType(err, TypeDataNotFound)
}

func NewDataNotFoundErr() error {
	err := errors.New("data not found")
	return serrors.Builder(err).
		Type(TypeDataNotFound).
		Build()
}

func NewDataNotFoundErrWithID(id string) error {
	return NewDataNotFoundErrWithKey("id", id)
}

func NewDataNotFoundErrWithKey(key, value string) error {
	err := NewDataNotFoundErr()
	return serrors.Builder(err).
		AddTagString(key, value).
		Build()
}
