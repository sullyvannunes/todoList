package main

import (
	"os"

	"github.com/sullyvannunes/todo-list/views"
)

func main() {
	// values := []struct {
	// 	Name string
	// }{
	// 	{"Sullyvan"},
	// 	{"Alice"},
	// 	{"Darwin"},
	// }
	// if err := views.Render(os.Stdout, "lists/index", values); err != nil {
	// 	panic(err)
	// }

	// t1 := template.Must(template.New("asdf").Parse(`Root: {{block "A" .}}This is the first template{{end}} {{block "B" .}}This is the second template{{end}}`))
	views.Render(os.Stdout, "", nil)
}

func Ensure(err error) {
	if err != nil {
		panic(err)
	}
}

// ao crar um novo template uma nova estrutura é isntanciada mas sem nenhum conteúdo
