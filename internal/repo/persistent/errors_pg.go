// Package persistent provides PostgreSQL-specific error handlers.
package persistent

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL error codes.
const (
	UniqueViolation     = "23505"
	ForeignKeyViolation = "23503"
	NotNullViolation    = "23502"
	CheckViolation      = "23514"
)

// IsUniqueViolation returns true if the error is a unique violation.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == UniqueViolation
}
