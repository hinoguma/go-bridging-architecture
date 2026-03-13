package errors

import (
	"context"

	serrors "github.com/hinoguma/go-structured-error"
)

func New(message string) error {
	return serrors.New(message)
}

func NewWithCtx(message string, ctx context.Context) error {
	return AddCtx(New(message), ctx)
}

func Wrap(err error, message string) error {
	return serrors.Wrap(err, message)
}

func Lift(err error) error {
	return serrors.Lift(err)
}

func LiftWithCtx(err error, ctx context.Context) error {
	return AddCtx(Lift(err), ctx)
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

func IsType(err error, targetType serrors.ErrorType) bool {
	return serrors.IsType(err, targetType)
}

func AddCtx(err error, ctx context.Context) error {
	rid := ctx.Value("requestId")
	if ridStr, ok := rid.(string); ok {
		return serrors.Builder(err).AddTagString("requestId", ridStr).Build()
	}
	return err
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
