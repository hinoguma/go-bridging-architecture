package crosscutting

type Stringer interface {
	String() string
}

type Stringable interface {
	FromString(s string)
}

func ConvertToStrings[T Stringer](items []T) []string {
	if items == nil {
		return make([]string, 0)
	}
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.String()
	}
	return result
}
func ConvertFromStrings[T any](items []string, convert func(s string) T) []T {
	if items == nil {
		return make([]T, 0)
	}
	result := make([]T, len(items))
	for i, item := range items {
		result[i] = convert(item)
	}
	return result
}
