package crosscutting

func Ptr[T any](value T) *T {
	if value == nil {
		return nil
	}
	return &value
}
