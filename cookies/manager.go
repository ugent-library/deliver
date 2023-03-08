package cookies

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
)

type cookie struct {
	data    any
	expires time.Time
}

type Manager struct {
	cookies map[string]cookie
}

func (m *Manager) Set(name string, data any, expires time.Time) {
	m.cookies[name] = cookie{data: data, expires: expires}
}

func (m *Manager) Append(name string, data any, expires time.Time) {
	var d []any
	if c, ok := m.cookies[name]; ok {
		if dd, ok := c.data.([]any); ok {
			d = dd
		}
	}
	d = append(d, data)
	m.cookies[name] = cookie{data: d, expires: expires}
}

func (m *Manager) Delete(name string) {
	m.Set(name, "", time.Now())
}

func (m *Manager) Write(w http.ResponseWriter) {
	for name, c := range m.cookies {
		var value string
		if v, ok := c.data.(string); ok {
			value = v
		} else {
			b, err := json.Marshal(c.data)
			if err != nil {
				return
			}
			value = base64.URLEncoding.EncodeToString(b)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    value,
			Expires:  c.expires,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

func Unmarshal(r *http.Request, name string, data any) error {
	if cookie, _ := r.Cookie(name); cookie != nil {
		b, err := base64.URLEncoding.DecodeString(cookie.Value)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, data); err != nil {
			return err
		}
	}
	return nil
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var sessionKey = contextKey("cookies")

func Jar(r *http.Request) *Manager {
	if v := r.Context().Value(sessionKey); v != nil {
		return v.(*Manager)
	}
	return nil
}

// futureproofing, not used yet
type options struct{}

type Option func(*options)

func Manage(opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := &Manager{cookies: make(map[string]cookie)}
			written := false

			w = httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(nextFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						if !written {
							m.Write(w)
							written = true
						}
						nextFunc(code)
					}
				},
				Write: func(nextFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(b []byte) (int, error) {
						if !written {
							m.Write(w)
							written = true
						}
						return nextFunc(b)
					}
				},
			})
			c := context.WithValue(r.Context(), sessionKey, m)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}
