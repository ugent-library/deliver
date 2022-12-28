package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// TODO add convenience Flash methods?
type Session interface {
	Get(string) any
	Set(string, any)
	HasKey(string) bool
	Keys() []string
	Delete(string)
	Save() error
}

type SugaredSession struct {
	session Session
	changed bool
}

func NewSugaredSession(s Session) *SugaredSession {
	return &SugaredSession{session: s}
}

func (s *SugaredSession) Get(k string) any {
	return s.session.Get(k)
}

func (s *SugaredSession) HasKey(k string) bool {
	return s.session.HasKey(k)
}

func (s *SugaredSession) Keys() []string {
	return s.session.Keys()
}

func (s *SugaredSession) Pop(k string) any {
	v := s.session.Get(k)
	s.Delete(k)
	return v
}

func (s *SugaredSession) Set(k string, v any) {
	s.session.Set(k, v)
	s.changed = true
}

func (s *SugaredSession) Append(k string, v any) {
	if s.session.HasKey(k) {
		vals := s.session.Get(k).([]any)
		s.Set(k, append(vals, v))
	} else {
		s.Set(k, []any{v})
	}
}

func (s *SugaredSession) Delete(k string) {
	if s.session.HasKey(k) {
		s.changed = true
	}
	s.session.Delete(k)
}

func (s *SugaredSession) Clear() {
	keys := s.session.Keys()
	if len(keys) > 0 {
		s.changed = true
	}
	for _, k := range keys {
		s.session.Delete(k)
	}
}

func (s *SugaredSession) Save() error {
	if s.changed {
		if err := s.session.Save(); err != nil {
			return err
		}
		s.changed = false
	}
	return nil
}

type gorillaSession struct {
	session *sessions.Session
	req     *http.Request
	res     http.ResponseWriter
}

func (s *gorillaSession) Get(k string) any {
	return s.session.Values[k]
}

func (s *gorillaSession) HasKey(k string) bool {
	_, ok := s.session.Values[k]
	return ok
}

func (s *gorillaSession) Keys() []string {
	keys := make([]string, len(s.session.Values))
	i := 0
	for k := range s.session.Values {
		keys[i] = k.(string)
		i++
	}
	return keys
}

func (s *gorillaSession) Set(k string, v any) {
	s.session.Values[k] = v
}

func (s *gorillaSession) Delete(k string) {
	delete(s.session.Values, k)
}

func (s *gorillaSession) Save() error {
	return s.session.Save(s.req, s.res)
}
