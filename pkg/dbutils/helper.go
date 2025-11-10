package dbutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/errors"
)

func EncodeCursor(objectID string) (res string) {
	return base64.StdEncoding.EncodeToString([]byte(objectID))
}

func DecodeCursor(cursor string) (res string, err error) {
	dst, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", errors.NewValidationError(fmt.Errorf("invalid cursor: %w", err))
	}
	res = string(dst)
	return res, nil
}

func IsDuplicatedData(err error) bool {
	if err == nil {
		return false
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}

	if pqErr.Code == "23505" {
		return true
	}

	return false
}

const (
	QueryDefaultLimit = 100
)

// Helper function to convert empty JSON strings to nil for PostgreSQL JSONB
func ToJSONB(v interface{}) json.RawMessage {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// IsNullJSON checks if a json.RawMessage represents a JSON null value
func IsNullJSON(data json.RawMessage) bool {
	if len(data) == 0 {
		return true
	}
	// Check if it's literally "null" JSON value
	trimmed := string(data)
	return trimmed == "null" || trimmed == ""
}
