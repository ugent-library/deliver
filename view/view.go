package view

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"text/template"
)

var (
	FS                = os.DirFS("templates")
	TemplateExtension = ".gohtml"
	ContentType       = "text/html"
	FuncMap           template.FuncMap
	Option            = "missingkey=error"
	bufPool           = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

type Data = map[string]any

type View struct {
	Template *template.Template
	status   int
}

func MustNew(tmpl string, files ...string) View {
	v, err := New(tmpl, files...)
	if err != nil {
		panic(err)
	}
	return v
}

func New(tmpl string, files ...string) (View, error) {
	tmpl = tmpl + TemplateExtension

	for i, f := range files {
		files[i] = f + TemplateExtension
	}

	t, err := template.New(tmpl).
		Option(Option).
		Funcs(FuncMap).
		ParseFS(FS, append(files, tmpl)...)
	if err != nil {
		return View{}, err
	}

	return View{Template: t}, nil
}

func (v View) Status(code int) View {
	v.status = code
	return v
}

func (v View) Render(w http.ResponseWriter, data any) error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", ContentType)
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if err := v.Template.Execute(buf, data); err != nil {
		return err
	}

	if v.status != 0 {
		w.WriteHeader(v.status)
	}

	_, err := io.Copy(w, buf)
	return err
}
