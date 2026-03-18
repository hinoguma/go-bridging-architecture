package crosscutting

import "fmt"

func StringsToJsonString(strs []string) string {
	s := "["
	for i, str := range strs {
		if i > 0 {
			s += ","
		}
		s += `"` + str + `"`
	}
	s += "]"
	return s
}

func NumbersToJsonString[
	T int | int8 | int16 | int32 | int64 | float32 | float64,
](nums []any) string {
	s := "["
	for i, num := range nums {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%d", num)
	}
	s += "]"
	return s
}
