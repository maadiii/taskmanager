package errors

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsConstraint(err error, consName string) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok &&
		consName == pgErr.ConstraintName {
		return true
	}

	return false
}

func UniqueViolation(err error, tableName, columnName string) error {
	msg := fmt.Sprintf("%s_%s already exists", tableName, columnName)

	return &Error{
		error: New(err.Error()),
		key:   makeErrorKey(msg),
		code:  alreadyExists,
	}
}
