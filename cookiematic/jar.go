package cookiematic

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type cookie struct {
	data    any
	expires time.Time
}

type Jar struct {
	w       http.ResponseWriter
	r       *http.Request
	written bool
	cookies map[string]cookie
}

func NewJar(w http.ResponseWriter, r *http.Request) *Jar {
	return &Jar{
		w:       w,
		r:       r,
		cookies: make(map[string]cookie),
	}
}

func (m *Jar) Get(name string) string {
	if c, _ := m.r.Cookie(name); c != nil {
		return c.Value
	}
	return ""
}

func (m *Jar) Unmarshal(name string, data any) error {
	if cookie, _ := m.r.Cookie(name); cookie != nil {
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

func (m *Jar) Set(name string, data any, expires time.Time) {
	m.cookies[name] = cookie{data: data, expires: expires}
}

func (m *Jar) Append(name string, data any, expires time.Time) {
	var d []any
	if c, ok := m.cookies[name]; ok {
		if dd, ok := c.data.([]any); ok {
			d = dd
		}
	}
	d = append(d, data)
	m.cookies[name] = cookie{data: d, expires: expires}
}

func (m *Jar) Delete(name string) {
	m.Set(name, "", time.Now())
}

func (m *Jar) Write() error {
	if m.written {
		return nil
	} else {
		m.written = true
	}

	for name, c := range m.cookies {
		var value string
		if v, ok := c.data.(string); ok {
			value = v
		} else {
			b, err := json.Marshal(c.data)
			if err != nil {
				return err
			}
			value = base64.URLEncoding.EncodeToString(b)
		}
		http.SetCookie(m.w, &http.Cookie{
			Name:     name,
			Value:    value,
			Expires:  c.expires,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	return nil
}
