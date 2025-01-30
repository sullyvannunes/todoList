package application

import (
	"io"
	"log"
	"net/http"
	"os"
)

type RendererFuncMap map[string]any

type Renderer interface {
	Render(io.Writer, string, any, RendererFuncMap) error
	RenderWithLayout(io.Writer, string, string, any, RendererFuncMap) error
}

type Application struct {
	view        Renderer
	httpHandler http.Handler
}

func New(view Renderer) *Application {
	app := &Application{}

	mux := http.NewServeMux()
	mux.HandleFunc("/lists", app.ListIndex())
	mux.HandleFunc("/actions", app.ActionIndex())
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServerFS(os.DirFS("assets"))))

	app.httpHandler = mux
	app.view = view

	return app
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.httpHandler.ServeHTTP(w, r)
}

type User struct {
	Name string
}

type RenderData struct {
	User
}

func (a *Application) ListIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := RenderData{User: User{"Sullyvan"}}

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if err := a.view.Render(w, "list_index.html", data, RendererFuncMap{}); err != nil {
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

		if err := a.view.Render(w, "action_index.html", nil, RendererFuncMap{
			"Title": func() string {
				return "Suas Ações"
			},
		}); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
