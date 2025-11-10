package dbutils

import (
	"database/sql"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func HandlePostgresError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return errors.ErrDataNotFound
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}

	switch pqErr.Code {
	case "23505":
		log.Error().
			Str("code", string(pqErr.Code)).
			Str("constraint", pqErr.Constraint).
			Str("detail", pqErr.Detail).
			Msg("Duplicate key violation")
		return errors.ErrDuplicatedData
	default:
		log.Error().
			Str("code", string(pqErr.Code)).
			Str("message", pqErr.Message).
			Str("detail", pqErr.Detail).
			Msg("Postgres error occurred")
	}

	return err
}
