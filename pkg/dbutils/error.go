package dbutils

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"

	apperrors "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
)

func HandlePostgresError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.ErrDataNotFound
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case "23505":
		log.Error().
			Str("code", string(pqErr.Code)).
			Str("constraint", pqErr.Constraint).
			Str("detail", pqErr.Detail).
			Msg("Duplicate key violation")
		return apperrors.ErrDuplicatedData
	default:
		log.Error().
			Str("code", string(pqErr.Code)).
			Str("message", pqErr.Message).
			Str("detail", pqErr.Detail).
			Msg("Postgres error occurred")
	}

	return err
}
