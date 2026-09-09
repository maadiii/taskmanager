package errors

import "fmt"

func NotFound(err error, what string) error {
	msg := fmt.Sprintf("%s not found", what)

	return &Error{
		error: New(err.Error()),
		key:   makeErrorKey(msg),
		code:  notFound,
	}
}

func Forbidden() error {
	msg := "forbidden"

	return &Error{
		error: New(msg),
		key:   makeErrorKey(msg),
		code:  forbidden,
	}
}
