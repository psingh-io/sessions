package sqlstore_test

import (
	"database/sql"
	"os"
	"sessions/store/sqlstore"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db *sql.DB
var store *sqlstore.SqlStore

type testSessionData struct {
	token uuid.UUID
	data  string
}

var testData testSessionData

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test")
	defer db.Close()

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	store = sqlstore.New(db)

	os.Exit(m.Run())
}

func beforeEach(t *testing.T, testData *testSessionData) func() {
	testData.token = uuid.New()
	testData.data = "{\"Token\": \"df25799c-5219-4c1a-b7b7-820250eddf66\", \"Values\": {\"hello\": \"message\"}}"

	err := store.Save(testData.token, testData.data, time.Now().Add(time.Minute*10).Unix())

	if err != nil {
		t.Fatal(err)
	}

	return func() {
		_, err := db.Exec("TRUNCATE TABLE sessions")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestMySqlStore_Save(t *testing.T) {
	testData = testSessionData{}
	defer beforeEach(t, &testData)()

	stmt := "SELECT data FROM sessions WHERE token = ?"
	row := db.QueryRow(stmt, testData.token)

	var data string
	err := row.Scan(&data)

	if err == nil {
		if testData.data != data {
			t.Error("Session data retieved from db does not match")
		} else {
			// t.Log("Success")
		}
	} else {
		t.Error(err)
	}
}

func TestMySqlStore_Get(t *testing.T) {
	testData = testSessionData{}
	defer beforeEach(t, &testData)()

	data, found, err := store.Get(testData.token)

	if err != nil {
		t.Error(err)
	} else if !found {
		t.Errorf("Token %s not found", testData.token)
	} else if testData.data != data {
		t.Error("Session data retieved from db does not match")
	} else {
		// t.Log("Success")
	}
}

func TestMySqlStore_Delete(t *testing.T) {
	testData = testSessionData{}
	defer beforeEach(t, &testData)()

	err := store.Delete(testData.token)

	data, found, err := store.Get(testData.token)
	if data != "" {
		t.Errorf("Expecting empty string, found %s", data)
	} else if found {
		t.Error("Deleted record found")
	} else if err != nil {
		t.Errorf("Error in getting deleted record")
	}

	stmt := "SELECT data FROM sessions WHERE token = ?"
	row := db.QueryRow(stmt, testData.token)

	err = row.Scan(&data)

	if err != sql.ErrNoRows {
		t.Errorf("Expecting ErrNoRows error, found %s", err)
	}
}

func TestMySqlStore_DeleteExpired(t *testing.T) {
	testData = testSessionData{}
	defer beforeEach(t, &testData)()

	expiry := time.Now().Unix()
	data := "{\"test\": \"value\"}"
	_, err := db.Exec("INSERT INTO sessions (token, data, expiry, createtime, updatetime) VALUES (?, ?, ?, ?, ?) ",
		uuid.New(), data, expiry-1, time.Now(), time.Now())
	if err != nil {
		t.Error(err)
	}
	_, _ = db.Exec("INSERT INTO sessions (token, data, expiry, createtime, updatetime) VALUES (?, ?, ?, ?, ?) ",
		uuid.New(), data, expiry-2, time.Now(), time.Now())
	_, _ = db.Exec("INSERT INTO sessions (token, data, expiry, createtime, updatetime) VALUES (?, ?, ?, ?, ?) ",
		uuid.New(), data, expiry-3, time.Now(), time.Now())

	stmt := "SELECT COUNT(*) FROM sessions"
	row := db.QueryRow(stmt)

	var count int64
	row.Scan(&count)

	if count != 4 {
		t.Errorf("Expecting 4 rows, found %d", count)
	}

	count, _ = store.DeleteExpired(time.Now().Unix())

	if count != 3 { // first record had expiry now + 10 mins and will not be deleted
		t.Errorf("Expecting 3 rows to be deleted, deleted %d rows", count)
	}
}
