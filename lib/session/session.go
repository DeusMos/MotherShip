package session

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Session struct {
	Session    string
	Permission string
}

func FromJSON(sessionJson []byte) (s *Session, err error) {
	err = ioutil.WriteFile("./session", sessionJson, 0644)
	if err != nil {
		return
	}
	err = json.NewDecoder(bytes.NewReader(sessionJson)).Decode(&s)
	return
}
func Load() (s *Session, err error) {
	sessionJson, err := ioutil.ReadFile("./session")
	if err != nil {
		return
	}
	err = json.NewDecoder(bytes.NewReader(sessionJson)).Decode(&s)
	return
}
