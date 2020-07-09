package store

import (
	"github.com/google/uuid"
)

type Store interface {
	Get(token uuid.UUID) (string, bool, error)
	Save(token uuid.UUID, data string, expiry int64) error
	GetBytes(token uuid.UUID) ([]byte, bool, error)
	SaveBytes(token uuid.UUID, data []byte, expiry int64) error
	Delete(token uuid.UUID) error
	DeleteExpired(expiry int64) (int64, error)
}
