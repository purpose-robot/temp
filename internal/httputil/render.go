package httputil

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
)

type templateData struct {
	Form  any
	Flash string
}

func NewTemplateData() *templateData {
	return &templateData{}
}

type Renderer struct {
	templateFS      fs.FS
	sharedTemplates *template.Template
}

func NewRenderer(templateFS fs.FS, sharedTemplateFiles ...string) (*Renderer, error) {
	funcs := template.FuncMap{}

	sharedTemplates, err := template.New("").Funcs(funcs).ParseFS(templateFS, sharedTemplateFiles...)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		templateFS:      templateFS,
		sharedTemplates: sharedTemplates,
	}, nil
}

func (t *Renderer) Render(w http.ResponseWriter, status int, data any, templateName string, additionalTemplateFiles ...string) error {
	st, err := t.sharedTemplates.Clone()
	if err != nil {
		return err
	}

	if len(additionalTemplateFiles) > 0 {
		st, err = st.ParseFS(t.templateFS, additionalTemplateFiles...)
		if err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)

	if err := st.ExecuteTemplate(buf, templateName, data); err != nil {
		return err
	}

	w.Header().Add("Vary", "HX-Request-Type")

	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)

	return nil
}
