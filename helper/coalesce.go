package helper

import "strings"

func Coalesce[T comparable](vals ...T) T {
	var zero T
	for _, v := range vals {
		switch any(v).(type) {
		case string:
			if strings.TrimSpace(any(v).(string)) != "" {
				return v
			}
		default:
			if v != zero {
				return v
			}
		}
	}
	return zero
}
