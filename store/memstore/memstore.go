package memstore

import (
	"github.com/google/uuid"
)

type MemStore struct {
	values map[uuid.UUID] interface{}
}

func New() *MemStore {
	return &MemStore {
		values: make(map[uuid.UUID]interface{}),
	}
}

func (m *MemStore) Get(token uuid.UUID) (string, bool, error) {
	val := m.values[token];
	if (val != nil) {
		return val.(string), true, nil
	} else {
		return "", false, nil
	}
}

func (m *MemStore) Save(token uuid.UUID, data string, expiry int64) error {
	m.values[token] = data
	return nil
}

func (m *MemStore) GetBytes(token uuid.UUID) ([]byte, bool, error) {
	val := m.values[token];
	if (val != nil) {
		return val.([]byte), true, nil
	} else {
		return nil, false, nil
	}
}

func (m *MemStore) SaveBytes(token uuid.UUID, data []byte, expiry int64) error {
	m.values[token] = data
	return nil
}

func (m *MemStore) Delete(token uuid.UUID) error {
	delete(m.values, token)
	return nil
}

func (m *MemStore) DeleteExpired(expiry int64) (int64, error) {
	return 0, nil
}