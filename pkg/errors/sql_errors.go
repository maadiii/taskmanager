package errors

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Pg(err *pgconn.PgError) error {
	switch err.Code {
	case pgerrcode.UniqueViolation:
		return uniqueViolation(err)
	case pgerrcode.CaseNotFound:
		return notfound(err)
	default:
		return err
	}
}

func uniqueViolation(err *pgconn.PgError) error {
	name := fmt.Sprintf("%s_%s", err.TableName, err.ColumnName)

	return &Error{
		error: errors.New(err.Error()),
		key:   makeErrorKey(name),
	}
}

func notfound(err *pgconn.PgError) error {
}
