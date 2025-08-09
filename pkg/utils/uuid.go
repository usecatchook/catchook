package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(id string) (pgtype.UUID, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	return pgtype.UUID{
		Bytes: userUUID,
		Valid: true,
	}, nil
}
