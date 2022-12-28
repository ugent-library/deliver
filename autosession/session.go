package autosession

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Session interface {
	HasKey(string) bool
	Keys() []string
	Get(string) any
	Pop(string) any
	Set(string, any)
	Append(string, any)
	Delete(string)
	Clear()
	Save() error
}

type SessionFunc func(http.ResponseWriter, *http.Request) Session

func GorillaSession(store sessions.Store, name string) SessionFunc {
	return func(w http.ResponseWriter, r *http.Request) Session {
		// TODO handle error
		s, _ := store.Get(r, name)
		return &gorillaSession{
			session: s,
			w:       w,
			r:       r,
		}
	}
}

type gorillaSession struct {
	session *sessions.Session
	w       http.ResponseWriter
	r       *http.Request
	changed bool
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

func (s *gorillaSession) Get(k string) any {
	return s.session.Values[k]
}

func (s *gorillaSession) Pop(k string) any {
	v := s.Get(k)
	s.Delete(k)
	return v
}

func (s *gorillaSession) Set(k string, v any) {
	s.session.Values[k] = v
	s.changed = true
}

func (s *gorillaSession) Append(k string, v any) {
	if s.HasKey(k) {
		vals := s.Get(k).([]any)
		s.Set(k, append(vals, v))
	} else {
		s.Set(k, []any{v})
	}
}

func (s *gorillaSession) Delete(k string) {
	if s.HasKey(k) {
		s.changed = true
	}
	delete(s.session.Values, k)
}

func (s *gorillaSession) Clear() {
	keys := s.Keys()
	if len(keys) > 0 {
		s.changed = true
	}
	for _, k := range keys {
		delete(s.session.Values, k)
	}
}

func (s *gorillaSession) Save() error {
	if s.changed {
		if err := s.session.Save(s.r, s.w); err != nil {
			return err
		}
		s.changed = false
	}
	return nil
}
