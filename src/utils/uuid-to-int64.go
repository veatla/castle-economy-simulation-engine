package utils

import (
	"encoding/binary"

	"github.com/google/uuid"
)

func UUIDToInt64(u uuid.UUID) int64 {
	return -int64(binary.BigEndian.Uint64(u[8:16]))
}
