package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func EntityAlreadyExists(name string) error {
	msg := fmt.Sprintf("%s already exists", name)

	return &Error{
		error: errors.New(msg),
		key:   makeErrorKey(msg),
	}
}

func EntityNotFound(name string) error {
	msg := fmt.Sprintf("%s not found", name)

	return &Error{
		error: errors.New(msg),
		key:   makeErrorKey(msg),
	}
}
