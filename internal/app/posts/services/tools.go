package services

import (
	"fmt"

	"github.com/google/uuid"
)

func StrToUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parsing uid: %s", err)
	}

	return parsedUUID, nil
}
