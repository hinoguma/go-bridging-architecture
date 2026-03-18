package errors

import serrors "github.com/hinoguma/go-structured-error"

func ToJsonString(err error) string {
	return serrors.ToJsonString(err)
}
