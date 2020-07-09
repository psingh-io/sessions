package sqlstore

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SqlStore struct {
	db *sql.DB
}

func New(db *sql.DB) *SqlStore {
	return &SqlStore{
		db: db,
	}
}

func (store *SqlStore) Get(token uuid.UUID) (string, bool, error) {
	var data string
	var stmt string

	stmt = "SELECT data FROM sessions WHERE token = ? AND expiry >= ?"

	row := store.db.QueryRow(stmt, token, time.Now().Unix())
	err := row.Scan(&data)
	if err == sql.ErrNoRows {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return data, true, nil
}

func (store *SqlStore) Save(token uuid.UUID, data string, expiry int64) error {
	_, err := store.db.Exec("INSERT INTO sessions (token, data, expiry, createtime, updatetime) VALUES (?, ?, ?, ?, ?) ",
		token, data, expiry, time.Now(), time.Now())

	return err
}

func (store *SqlStore) GetBytes(token uuid.UUID) ([]byte, bool, error) {
	//var b []byte
	//var stmt string
	//
	//if compareVersion("5.6.4", m.version) >= 0 {
	//	stmt = "SELECT data FROM sessions WHERE token = ? AND UTC_TIMESTAMP(6) < expiry"
	//} else {
	//	stmt = "SELECT data FROM sessions WHERE token = ? AND UTC_TIMESTAMP < expiry"
	//}
	//
	//row := m.DB.QueryRow(stmt, token)
	//err := row.Scan(&b)
	//if err == sql.ErrNoRows {
	//	return nil, false, nil
	//} else if err != nil {
	//	return nil, false, err
	//}
	//return b, true, nil
	return nil, false, nil
}

func (store *SqlStore) SaveBytes(token uuid.UUID, data []byte, expiry int64) error {
	return nil
}

func (store *SqlStore) Delete(token uuid.UUID) error {
	_, err := store.db.Exec("DELETE FROM sessions where token = ?",
		token)
	return err
}

func (store *SqlStore) DeleteExpired(expiry int64) (int64, error) {
	res, err := store.db.Exec("DELETE FROM sessions where expiry < ?",
		expiry)

	if err == nil {
		return res.RowsAffected()
	} else {
		return -1, err
	}
}
