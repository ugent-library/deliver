package autosession

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionProvider func(http.ResponseWriter, *http.Request) (*Session, error)

func GorillaSessions(store sessions.Store, name string) SessionProvider {
	return func(w http.ResponseWriter, r *http.Request) (*Session, error) {
		s, err := store.Get(r, name)
		if err != nil {
			return nil, err
		}
		return &Session{
			session: s,
			w:       w,
			r:       r,
		}, nil
	}
}

type Session struct {
	session *sessions.Session
	w       http.ResponseWriter
	r       *http.Request
	changed bool
}

func (s *Session) HasKey(k string) bool {
	_, ok := s.session.Values[k]
	return ok
}

func (s *Session) Keys() []string {
	keys := make([]string, len(s.session.Values))
	i := 0
	for k := range s.session.Values {
		keys[i] = k.(string)
		i++
	}
	return keys
}

func (s *Session) Get(k string) any {
	return s.session.Values[k]
}

func (s *Session) Pop(k string) any {
	v := s.Get(k)
	s.Delete(k)
	return v
}

func (s *Session) Set(k string, v any) {
	s.session.Values[k] = v
	s.changed = true
}

func (s *Session) Append(k string, v any) {
	if s.HasKey(k) {
		vals := s.Get(k).([]any)
		s.Set(k, append(vals, v))
	} else {
		s.Set(k, []any{v})
	}
}

func (s *Session) Delete(k string) {
	if s.HasKey(k) {
		s.changed = true
	}
	delete(s.session.Values, k)
}

func (s *Session) Clear() {
	keys := s.Keys()
	if len(keys) > 0 {
		s.changed = true
	}
	for _, k := range keys {
		delete(s.session.Values, k)
	}
}

func (s *Session) Save() error {
	if s.changed {
		if err := s.session.Save(s.r, s.w); err != nil {
			return err
		}
		s.changed = false
	}
	return nil
}
