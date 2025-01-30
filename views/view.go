package views

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sullyvannunes/todo-list/application"
)

type View struct {
	viewsMap          map[string]*template.Template
	productionMode    bool
	namespace         string
	defaultLayoutName string
	componentDir      string
}

func New() *View {
	return &View{
		viewsMap:          make(map[string]*template.Template),
		productionMode:    false,
		namespace:         "views",
		defaultLayoutName: "application.html",
		componentDir:      "components",
	}
}

func (v *View) Render(w io.Writer, tmplName string, data any, funcMap application.RendererFuncMap) error {
	return v.RenderWithLayout(w, v.defaultLayoutName, tmplName, data, funcMap)
}

func (v *View) execute(tmpl *template.Template, w io.Writer, data any) error {
	return tmpl.Execute(w, data)
}

func (v *View) RenderWithLayout(w io.Writer, layoutName, tmplName string, data any, funcMap application.RendererFuncMap) error {
	layout := template.New(layoutName).Funcs(v.viewFuncs()).Funcs(template.FuncMap(funcMap))
	_, err := layout.ParseFiles(filepath.Join(v.namespace, layoutName))
	if err != nil {
		return fmt.Errorf("parse layout %w", err)
	}

	if err = v.parseTemplate(layout, filepath.Join(v.namespace, tmplName)); err != nil {
		return err
	}

	fsys := os.DirFS(filepath.Join(v.namespace, v.componentDir))
	if _, err = layout.ParseFS(fsys, "*.html"); err != nil {
		return fmt.Errorf("parse components: %w", err)
	}

	if err = layout.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func (v *View) parseTemplate(base *template.Template, templatePath string) error {
	if _, err := base.ParseFiles(templatePath); err != nil {
		return fmt.Errorf("parse template %s: %w", templatePath, err)
	}

	return nil
}

func (v *View) viewFuncs() template.FuncMap {
	return template.FuncMap{}
}
