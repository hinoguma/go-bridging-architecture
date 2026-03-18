package errors

import (
	"context"

	serrors "github.com/hinoguma/go-structured-error"
)

// New() / NewWithCtx() make a new structured error with the message and stack trace.
func New(message string) error {
	return serrors.New(message)
}

func NewWithCtx(message string, ctx context.Context) error {
	return AddCtx(New(message), ctx)
}

/**
 * Wrap() / WrapWithCtx() / Lift() / LiftWithCtx()
 * These functions convert error to structured error.
 * Structured error has stack trace and can have tags and sub errors.
 * If the error is already a structured error, it returns the same error.
 */

func Wrap(err error, message string) error {
	return serrors.Wrap(err, message)
}

func WrapWithCtx(err error, message string, ctx context.Context) error {
	return AddCtx(Wrap(err, message), ctx)
}

func Lift(err error) error {
	return serrors.Lift(err)
}

func LiftWithCtx(err error, ctx context.Context) error {
	return AddCtx(Lift(err), ctx)
}

type ErrorType serrors.ErrorType

func (value ErrorType) SErrorType() serrors.ErrorType {
	return serrors.ErrorType(value)
}

func Is(err, target error) bool {
	return serrors.Is(err, target)
}

func As(err error, target any) bool {
	return serrors.As(err, target)
}

func Unwrap(err error) error {
	return serrors.Unwrap(err)
}

func IsType(err error, targetType ErrorType) bool {
	return serrors.IsType(err, targetType.SErrorType())
}

// AddCtx() adds RequestID from context to error as a tag if it exists in the context.
func AddCtx(err error, ctx context.Context) error {
	requestIDKey := GetContextRequestIDKey()
	rid := ctx.Value(requestIDKey)
	if ridStr, ok := rid.(string); ok {
		return serrors.Builder(err).AddTagString(requestIDKey, ridStr).Build()
	}
	return err
}

func AddType(err error, errType ErrorType) error {
	return serrors.Builder(err).Type(errType.SErrorType()).Build()
}

func AddTagString(err error, key, value string) error {
	return serrors.Builder(err).AddTagString(key, value).Build()
}

func AddSubErr(err error, adds ...error) error {
	return serrors.ToStructured(err).AddSubError(adds...)
}

type HasError struct {
	Err error
}

func (model HasError) IsErr() bool {
	return model.Err != nil
}

func DeferFuncAsSubErr(err error, deferFunc func() error) error {
	if deferErr := deferFunc(); deferErr != nil {
		if err == nil {
			return Lift(deferErr)
		}
		serr := serrors.ToStructured(err)
		return serr.AddSubError(deferErr)
	}
	return err
}
