package errors

import serrors "github.com/hinoguma/go-structured-error"

const TypeInvalidPassword serrors.ErrorType = "INVALID_PASSWORD"

func IsInvalidPasswordError(err error) bool {
	return serrors.IsType(err, TypeInvalidPassword)
}

func NewInvalidPasswordErrFromErr(err error) error {
	err = serrors.Lift(err)
	return serrors.Builder(err).
		Type(TypeInvalidPassword).
		Build()
}
