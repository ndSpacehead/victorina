package server

import (
	"bytes"
	"text/template"
	"net/http"
)

type templates struct {
	t *template.Template
}

func newTemplates() (*templates, error) {
	tmpl, err := template.ParseFS(indexHTML, "template/*.html")
	if err != nil {
		return nil, err
	}
	return &templates{
		t: tmpl,
	}, nil
}

func (t *templates) render(w http.ResponseWriter, name string, data any) error {
	buf := new(bytes.Buffer)
	defer buf.Reset()
	if err := t.t.ExecuteTemplate(buf, name, data); err != nil {
		return err
	}
	_, err := w.Write(buf.Bytes())
	return err
}
