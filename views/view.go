package views

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

type View struct {
	viewsMap       map[string]*template.Template
	productionMode bool
}

func New() *View {
	return &View{
		viewsMap:       make(map[string]*template.Template),
		productionMode: false,
	}
}

func (v *View) mergeTemplates(dst, src *template.Template) error {
	for _, t := range src.Templates() {
		if _, err := dst.AddParseTree(t.Name(), t.Tree); err != nil {
			return fmt.Errorf("add parse tree %s: %w", t.Name(), err)
		}
	}

	return nil
}

func (v *View) addTemplates(dst *template.Template, pattern string) error {
	fsys := os.DirFS("views")
	filenames, err := fs.Glob(fsys, pattern)
	if err != nil {
		return err
	}

	if len(filenames) == 0 {
		return nil
	}

	src, err := dst.ParseFS(fsys, pattern)
	if err != nil {
		return err
	}

	err = v.mergeTemplates(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func (v *View) Render(w io.Writer, tmplName string, data any) error {
	return v.RenderWithLayout(w, "application.html", tmplName, data)
}

func (v *View) Execute(tmpl *template.Template, w io.Writer, data any) error {
	return tmpl.Execute(w, data)
}

func (v *View) RenderWithLayout(w io.Writer, layoutName, tmplName string, data any) error {
	if v.productionMode {
		t, ok := v.viewsMap[tmplName]
		if ok {
			return v.Execute(t, w, data)
		}
	}

	layout, err := template.New("Root").Funcs(v.ViewFuncs()).ParseFiles("views/" + layoutName)
	if err != nil {
		return fmt.Errorf("parse layout: %w", err)
	}

	_, err = layout.ParseFiles("views/" + tmplName + ".html")
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplName, err)
	}

	fsys := os.DirFS("views")
	filenames, err := fs.Glob(fsys, "components/*.html")
	if err != nil {
		return err
	}

	if len(filenames) != 0 {
		_, err = layout.ParseFS(os.DirFS("views"), "components/*.html")
		if err != nil {
			return fmt.Errorf("parse components: %w", err)
		}
	}

	actionDir := strings.Split(tmplName, "/")[0]
	filenames, err = fs.Glob(fsys, actionDir+"/_*.html")
	if err != nil {
		return err
	}

	if len(filenames) != 0 {
		_, err = layout.ParseFS(os.DirFS("views"), actionDir+"/_*.html")
		if err != nil {
			return fmt.Errorf("parse components: %w", err)
		}
	}

	layoutTmpl := layout.Lookup(layoutName)

	if err = layoutTmpl.Execute(w, data); err != nil {
		return nil
	}

	v.viewsMap[tmplName] = layoutTmpl
	return nil
}

func (v *View) ViewFuncs() template.FuncMap {
	return template.FuncMap{
		"Title": func() string {
			return "Generated title"
		},
	}
}
