package server

import (
	"bytes"
	"math"
	"net/http"
	"text/template"
)

type templates struct {
	t *template.Template
}

func newTemplates() (*templates, error) {
	tmpl, err := template.New("vic").Funcs(template.FuncMap{
		"declensionScores": declensionScores,
	}).ParseFS(indexHTML, "template/*.html")
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

func declensionScores(score int) string {
	if score < 0 { // showing off
		if score == math.MinInt {
			return "баллов"
		}
		score = ^score + 1
	}
	score = score % 100
	if score > 4 && score < 21 {
		return "баллов"
	}
	switch score % 10 {
	case 1:
		return "балл"
	case 2, 3, 4:
		return "балла"
	default:
		return "баллов"
	}
}
