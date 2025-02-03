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
	mux.HandleFunc("GET /lists", app.ListIndex())
	mux.HandleFunc("POST /lists", app.CreateList())
	mux.HandleFunc("GET /actions", app.ActionIndex())
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

type List struct {
	Id   int64
	Name string
}

type RenderData struct {
	User
	Lists  []List
	Errors map[string][]string
}

func NewRenderData() RenderData {
	return RenderData{
		Errors: make(map[string][]string),
	}
}

var lists = make([]List, 0)

func (a *Application) ListIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := RenderData{User: User{"Sullyvan"}, Lists: []List{}}

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

func (a *Application) CreateList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := NewRenderData()
		data.User = User{"Sullyvan"}
		listName := r.FormValue("lists[name]")

		if listName == "" {
			w.WriteHeader(http.StatusBadRequest)
			data.Errors["name"] = []string{"O nome da lista não pode ser vazio"}
			a.view.Render(w, "list_index.html", data, RendererFuncMap{})
			return
		}

		lists = append(lists, List{Name: listName})
		data.Lists = lists
		a.view.Render(w, "list_index.html", data, RendererFuncMap{})

		w.Header().Set("Location", "/lists")
		w.WriteHeader(http.StatusFound)
	}
}
