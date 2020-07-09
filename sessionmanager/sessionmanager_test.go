package sessionmanager_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"

	"sessions/session"
)

// var handler http.Handler
var key string = "message"
var value = "Hello from a session!"

// var value string
var mux *http.ServeMux

func TestMain(m *testing.M) {
	mux = http.NewServeMux()

	mux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(*session.Session)
		session.Put(key, value)
	})

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(*session.Session)
		msg, _ := session.GetString("message")
		io.WriteString(w, msg)
	})
	os.Exit(m.Run())
}

func TestSessionManager_Put(t *testing.T) {

	sm := session.NewSessionManager()
	handler := sm.Handler(mux)

	r, err := http.NewRequest("GET", "/put", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	cookie := getSessionCookie(resp)

	if cookie == nil {
		t.Error("Session cookie not found")
	}

	uuid, _ := uuid.Parse(cookie.Value)
	s, found, _ := sm.Store.Get(uuid)

	if s == "" || !found {
		t.Errorf("Token %s not found", uuid.String())
	}

	var session session.Session
	json.Unmarshal([]byte(s), &session)
	v, _ := session.GetString(key)

	if v != value {
		t.Errorf("Stored value %s does not match session value %s", v, value)
	}

}

func TestSessionManager_Get(t *testing.T) {

	sm := session.NewSessionManager()
	handler := sm.Handler(mux)

	r, err := http.NewRequest("GET", "/put", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	cookie := getSessionCookie(resp)

	if cookie == nil {
		t.Error("Session cookie not found")
	}

	r, err = http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.AddCookie(cookie)

	handler.ServeHTTP(w, r)

	resp = w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	bodyString := w.Body.String()
	if bodyString != value {
		t.Errorf("Stored value \"%s\" does not match fetched value \"%s\"", bodyString, value)
	}
}

func getSessionCookie(resp *http.Response) *http.Cookie {
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "Session" {
			return cookie
		}
	}

	return nil
}
