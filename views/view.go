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
	return v.RenderWithLayout(w, "views/application.html", tmplName, data)
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

	file, _ := os.Open("views/" + tmplName + ".html")
	b, _ := io.ReadAll(file)
	tmpl, err := template.New("views/" + tmplName + ".html").Funcs(template.FuncMap{"Xablau": v.Xablau}).Parse(string(b))
	if err != nil {
		return err
	}

	if err = v.addTemplates(tmpl, "components/*.html"); err != nil {
		return fmt.Errorf("add template %s: %w", "components/*.html", err)
	}

	prefix := strings.Split(tmplName, "/")[0]
	if err = v.addTemplates(tmpl, prefix+"/_*.html"); err != nil {
		return fmt.Errorf("add template %s: %w", prefix+"/_*.html", err)
	}

	layout, err := template.ParseFiles(layoutName)
	if err != nil {
		return err
	}

	if err = v.mergeTemplates(layout, tmpl); err != nil {
		return err
	}

	layout.Funcs(template.FuncMap{
		"Xablau": v.Xablau,
	})

	if err = v.Execute(layout, w, data); err != nil {
		return err
	}

	v.viewsMap[tmplName] = layout

	return nil
}

func (v *View) Xablau(data any) string {
	fmt.Printf("%s %+v\n", "this is xablau", data)
	return "asdf"
}
