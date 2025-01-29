package main

import (
	"io"
	"log"
	"net/http"

	"github.com/sullyvannunes/todo-list/views"
)

type Renderer interface {
	Render(io.Writer, string, any) error
	RenderWithLayout(io.Writer, string, string, any) error
}

type Application struct {
	view        Renderer
	httpHandler http.Handler
}

func NewApplication(view Renderer) *Application {
	app := &Application{}

	mux := http.NewServeMux()
	mux.HandleFunc("/lists", app.ListIndex())
	mux.HandleFunc("/actions", app.ActionIndex())

	app.httpHandler = mux
	app.view = view

	return app
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.httpHandler.ServeHTTP(w, r)
}

func (a *Application) ListIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if err := a.view.Render(w, "lists/index", []int{1, 2, 3, 4, 5}); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (a *Application) ActionIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if err := a.view.Render(w, "actions/index", nil); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	view := views.New(views.ViewConfig{
		Namespace:             "views",
		ComponentDir:          "components",
		IsProductionMode:      false,
		LayoutDefault:         "application",
		TemplateFileExtension: ".html",
	})

	app := NewApplication(view)

	log.Println("running on http://localhost:3030")
	http.ListenAndServe(":3030", app)
}
