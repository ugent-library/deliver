package view

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

var DefaultConfig = Config{
	FS:                 os.DirFS("templates"),
	TemplateExtension:  ".gohtml",
	Option:             "missingkey=error",
	DefaultContentType: "text/html",
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type Config struct {
	FS                 fs.FS
	TemplateExtension  string
	Funcs              template.FuncMap
	Option             string
	DefaultContentType string
}

func (c Config) NewView(layout string, files ...string) (View, error) {
	files = append(files, layout)
	for i, f := range files {
		files[i] = f + c.TemplateExtension
	}

	tmpl, err := template.New("").
		Option(c.Option).
		Funcs(c.Funcs).
		ParseFS(c.FS, files...)
	if err != nil {
		return View{}, err
	}

	return View{
		layout: filepath.Base(layout),
		tmpl:   tmpl,
	}, nil
}

type View struct {
	tmpl        *template.Template
	layout      string
	contentType string
	status      int
}

func MustNew(layout string, files ...string) View {
	v, err := DefaultConfig.NewView(layout, files...)
	if err != nil {
		panic(err)
	}
	return v
}

func New(layout string, files ...string) (View, error) {
	return DefaultConfig.NewView(layout, files...)
}

func (v View) Status(s int) View {
	v.status = s
	return v
}

func (v View) ContentType(ct string) View {
	v.contentType = ct
	return v
}

func (v View) Render(w http.ResponseWriter, data any) error {
	header := w.Header()
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", v.contentType)
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if err := v.tmpl.ExecuteTemplate(buf, v.layout, data); err != nil {
		return err
	}

	if v.status != 0 {
		w.WriteHeader(v.status)
	}

	_, err := io.Copy(w, buf)
	return err
}
