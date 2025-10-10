package helper

import "github.com/jinzhu/copier"

// Convert copies matching fields from src into a new T.
// If transform is provided, it will be called after copying.
func Convert[T any](src any, transform ...func(*T) error) (T, error) {
	var dst T

	err := copier.CopyWithOption(&dst, src, copier.Option{
		DeepCopy:    true,
		IgnoreEmpty: false, // set true if you want to skip zero-value fields
	})
	if err != nil {
		return dst, err
	}

	if len(transform) > 0 && transform[0] != nil {
		if err := transform[0](&dst); err != nil {
			return dst, err
		}
	}
	return dst, nil
}

