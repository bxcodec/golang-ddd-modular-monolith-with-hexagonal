package dbutils

import (
	"database/sql"
	"log"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
	"github.com/lib/pq"
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
		log.Printf("ERROR: Duplicate key violation: %v", err)
		return errors.ErrDuplicatedData
	default:
		log.Printf("ERROR: Postgres error occurred: %v", err)
	}

	return err
}
