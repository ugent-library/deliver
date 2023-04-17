package crumb

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type responseCookie struct {
	data    any
	expires time.Time
}

type CookieJar struct {
	cookies         []*http.Cookie
	responseCookies map[string]responseCookie
}

func New(cookies []*http.Cookie) *CookieJar {
	return &CookieJar{
		cookies: cookies,
	}
}

func (m *CookieJar) Get(name string) string {
	for _, c := range m.cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func (m *CookieJar) Unmarshal(name string, data any) error {
	for _, c := range m.cookies {
		if c.Name == name {
			b, err := base64.URLEncoding.DecodeString(c.Value)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(b, data); err != nil {
				return err
			}
		}
	}
	return nil
}

func (j *CookieJar) Set(name string, data any, expires time.Time) {
	if j.responseCookies == nil {
		j.responseCookies = make(map[string]responseCookie)
	}
	j.responseCookies[name] = responseCookie{data: data, expires: expires}
}

func (j *CookieJar) Append(name string, data any, expires time.Time) {
	var d []any
	if c, ok := j.responseCookies[name]; ok {
		if dd, ok := c.data.([]any); ok {
			d = dd
		}
	}
	d = append(d, data)
	j.Set(name, d, expires)
}

func (j *CookieJar) Delete(name string) {
	j.Set(name, "", time.Now())
}

func (j *CookieJar) Write(w http.ResponseWriter) error {
	for name, c := range j.responseCookies {
		var val string
		switch v := c.data.(type) {
		case string:
			val = v
		case []byte:
			val = string(v)
		default:
			b, err := json.Marshal(c.data)
			if err != nil {
				return err
			}
			val = base64.URLEncoding.EncodeToString(b)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    val,
			Expires:  c.expires,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	return nil
}
