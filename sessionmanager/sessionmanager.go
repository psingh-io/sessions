package sessionmanager

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"sessions/session"
	"sessions/store"
	"sessions/store/memstore"
)

type list struct {
	Json   string
	Binary string
}

var EncodingTypes = &list{
	Json:   "Json",
	Binary: "Binary",
}

type CookieOptions struct {
	Name     string
	Domain   string
	HttpOnly bool
	Path     string
	Persist  bool
	Secure   bool
	SameSite http.SameSite
}

type SessionManager struct {
	IdleTimeout          time.Duration
	Encoding             string
	DirtyReadFromSession bool
	Store                store.Store
	CookieOptions        *CookieOptions
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		IdleTimeout:   10 * time.Minute,
		Encoding:      EncodingTypes.Json,
		Store:         memstore.New(),
		CookieOptions: defaultCookieOptions(),
	}
}

func (s *SessionManager) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var strToken string
		cookie, err := r.Cookie(s.CookieOptions.Name)
		if err == nil {
			strToken = cookie.Value
		} else {
			strToken = ""
		}

		var sn *session.Session
		if strToken != "" {
			token, _ := uuid.Parse(strToken)
			sn, _ = s.getSession(token)
		} else {
			sn = session.New()
		}

		c := context.WithValue(r.Context(), "session", sn)
		r2 := r.WithContext(c)

		next.ServeHTTP(w, r2)

		s.saveSession(sn)
		s.addSessionCookie(w, sn.Token)

	})
}

func (s *SessionManager) getSession(token uuid.UUID) (*session.Session, error) {
	var sess session.Session
	if s.Encoding == EncodingTypes.Binary {
		data, found, err := s.Store.GetBytes(token)
		if found && data != nil && err == nil {
			r := bytes.NewReader(data)
			err = gob.NewDecoder(r).Decode(&sess)

			if err == nil {
				return &sess, nil
			} else {
				return nil, err
			}
		} else {
			return session.New(), nil
		}
	} else {
		data, found, err := s.Store.Get(token)
		if found && data != "" && err == nil {

			err = json.Unmarshal([]byte(data), &sess)

			if err == nil {
				return &sess, nil
			} else {
				return nil, err
			}
		} else {
			return session.New(), err
		}
	}
}

type serializableSession struct {
	Token  string
	Values map[string]interface{}
}

func (s *SessionManager) saveSession(session *session.Session) error {
	if s.Encoding == EncodingTypes.Binary {
		var b bytes.Buffer
		err := gob.NewEncoder(&b).Encode(&session)
		if err == nil {
			err = s.Store.SaveBytes(session.Token, b.Bytes(), time.Now().Add(s.IdleTimeout).Unix())
			return err
		} else {
			return err
		}
	} else {
		data, err := json.Marshal(session)
		if err == nil {
			err = s.Store.Save(session.Token, string(data), time.Now().Add(s.IdleTimeout).Unix())
			return err
		} else {
			return err
		}
	}
}

func (s *SessionManager) addSessionCookie(w http.ResponseWriter, token uuid.UUID) {
	co := s.CookieOptions
	cookie := http.Cookie{
		Name:     co.Name,
		Value:    token.String(),
		Domain:   co.Domain,
		HttpOnly: co.HttpOnly,
		Path:     co.Path,
		Secure:   co.Secure,
		SameSite: co.SameSite,
	}

	if s.IdleTimeout != 0 {
		cookie.Expires = time.Now().Add(s.IdleTimeout)
		cookie.MaxAge = int(s.IdleTimeout / time.Second)
	} else {
		// Session cookie, default
	}

	http.SetCookie(w, &cookie)
}

func defaultCookieOptions() *CookieOptions {
	return &CookieOptions{
		Name:     "Session",
		Domain:   "",
		HttpOnly: true,
		Path:     "/",
		Persist:  false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
