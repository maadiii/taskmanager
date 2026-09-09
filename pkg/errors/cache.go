package errors

import "fmt"

func Cache(err error, operation string) error {
	return Wrap(fmt.Errorf("cache %s: %w", operation, err))
}
