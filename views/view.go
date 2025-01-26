package views

import (
	"embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed lists/*.html components/*.html layouts/*.html
var fs embed.FS

var (
	appTmpl = template.Must(template.ParseFiles("views/layouts/application.html"))
)

func Render(w io.Writer, tmplName string, data any) error {
	tmpl := template.Must(template.New("views").ParseFS(fs, "layouts/application.html", "lists/*.html"))

	for _, t := range tmpl.Templates() {
		fmt.Println(t.Name())
	}

	Ensure(tmpl.ExecuteTemplate(w, "application.html", data))

	return nil
}

func Ensure(err error) {
	if err != nil {
		panic(err)
	}
}
