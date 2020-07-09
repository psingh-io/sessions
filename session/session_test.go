package session_test

import (
	"sessions/session"
	"testing"
	"time"
)

func TestSession_GetString(t *testing.T) {
	sn := session.New()
	key := "Key"
	value := "Value"

	sn.Values[key] = value

	v, err := sn.GetString(key)

	if err != nil {
		t.Errorf("Error in getting string value for key %s", key)
	}

	if v != value {
		t.Errorf("String value mismatch. Expected \"%s\", Returned \"%s\"", value, v)
	}

	sn.Values[key] = 33
	v, err = sn.GetString(key)

	if err == nil {
		t.Error("Expected error, but no error returned")
	}
}

func TestSession_GetBool(t *testing.T) {
	sn := session.New()
	key := "Key"
	value := true

	sn.Values[key] = value

	v, err := sn.GetBool(key)

	if err != nil {
		t.Errorf("Error in getting bool value for key %s", key)
	}

	if v != value {
		t.Errorf("String value mismatch. Expected \"%v\", Returned \"%v\"", value, v)
	}

	sn.Values[key] = "true"
	v, err = sn.GetBool(key)

	if err != nil {
		t.Error(err)
	}

	sn.Values[key] = "false"
	v, err = sn.GetBool(key)

	if err != nil {
		t.Error(err)
	}

	sn.Values[key] = 33
	v, err = sn.GetBool(key)

	if err == nil {
		t.Error("Expected error, but no error returned")
	}
}

func TestSession_GetInt(t *testing.T) {
	sn := session.New()
	key := "Key"
	value := 500

	sn.Values[key] = value

	v, err := sn.GetInt(key)

	if err != nil {
		t.Errorf("Error in getting int value for key %s", key)
	}

	if v != value {
		t.Errorf("String value mismatch. Expected \"%v\", Returned \"%v\"", value, v)
	}

	sn.Values[key] = "33"
	v, err = sn.GetInt(key)

	if err != nil {
		t.Error("Expected no error, but error returned")
	}

	if v != 33 {
		t.Errorf("Expected 33, returned %d", v)
	}
}

func TestSession_GetFloat(t *testing.T) {
	sn := session.New()
	key := "Key"
	value := 500.33

	sn.Values[key] = value

	v, err := sn.GetFloat(key)

	if err != nil {
		t.Errorf("Error in getting float value for key %s", key)
	}

	if v != value {
		t.Errorf("String value mismatch. Expected \"%v\", Returned \"%v\"", value, v)
	}

	sn.Values[key] = "33.33"
	v, err = sn.GetFloat(key)

	if err != nil {
		t.Error("Expected no error, but error returned")
	}

	if v != 33.33 {
		t.Errorf("Expected 33.33, returned %f", v)
	}
}

func TestSession_GetTime(t *testing.T) {
	sn := session.New()
	key := "Key"
	value := time.Now()

	sn.Values[key] = value

	v, err := sn.GetTime(key)

	if err != nil {
		t.Errorf("Error in getting time value for key %s", key)
	}

	if v != value {
		t.Errorf("String value mismatch. Expected \"%v\", Returned \"%v\"", value, v)
	}

	s := value.Format(time.RFC3339)
	sn.Values[key] = s

	v, err = sn.GetTime(key)

	if err != nil {
		t.Error("Expected no error, but error returned")
	}

	time1, _ := time.Parse(time.RFC3339, s)
	if time1 != v {
		t.Errorf("Expected %v, returned %v", time1, v)
	}
}
