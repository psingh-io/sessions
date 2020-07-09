package session

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token     uuid.UUID
	Values    map[string]interface{}
	dirtyRead bool
	lock      sync.Mutex
}

func New() *Session {
	uuid, _ := uuid.NewRandom()
	return &Session{
		Token:  uuid,
		Values: make(map[string]interface{}),
	}
}

func NewWithDirtyRead(dirtyRead bool) *Session {
	uuid, _ := uuid.NewRandom()
	return &Session{
		Token:     uuid,
		Values:    make(map[string]interface{}),
		dirtyRead: dirtyRead,
	}
}

func NewWithValues(token uuid.UUID, values map[string]interface{}) *Session {
	if values != nil {
		return &Session{
			Token:  token,
			Values: values,
		}
	} else {
		return New()
	}
}

func (s *Session) Put(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Values[key] = value
}

func (s *Session) GetString(key string) (string, error) {
	if !s.dirtyRead {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	val := s.Values[key]
	str, ok := val.(string)

	if ok {
		return str, nil
	} else {
		return str, fmt.Errorf("%v is not a valid String", val)
	}
}

func (s *Session) GetBool(key string) (bool, error) {
	if !s.dirtyRead {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	val := s.Values[key]
	b, ok := val.(bool)

	if ok {
		return b, nil
	} else {
		str, ok := val.(string)
		if ok {
			if strings.EqualFold(str, "true") {
				return true, nil
			} else if strings.EqualFold(str, "false") {
				return false, nil
			}
		}
		return b, fmt.Errorf("%v is not a valid Bool", val)
	}
}

func (s *Session) GetInt(key string) (int, error) {
	if !s.dirtyRead {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	val := s.Values[key]
	i, ok := val.(int)

	if ok {
		return i, nil
	} else {
		str, ok := val.(string)
		if ok {
			i, err := strconv.Atoi(str)
			if err == nil {
				return i, nil
			}
		}
		return i, fmt.Errorf("%v is not a valid Integer", val)
	}
}

func (s *Session) GetFloat(key string) (float64, error) {
	if !s.dirtyRead {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	val := s.Values[key]
	f, ok := val.(float64)

	if ok {
		return f, nil
	} else {
		str, ok := val.(string)
		if ok {
			i, err := strconv.ParseFloat(str, 64)
			if err == nil {
				return i, nil
			}
		}
		return f, fmt.Errorf("%v is not a valid Float", val)
	}
}

func (s *Session) GetTime(key string) (time.Time, error) {
	if !s.dirtyRead {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	val := s.Values[key]
	t, ok := val.(time.Time)

	if ok {
		return t, nil
	} else {
		str, ok := val.(string)
		if ok {
			time1, err := time.Parse(time.RFC3339, str)
			if err == nil {
				return time1, nil
			}
		}

		return t, fmt.Errorf("%v is not a valid Time", val)
	}
}
