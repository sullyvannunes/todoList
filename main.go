package main

import (
	"log"
	"net/http"

	"github.com/sullyvannunes/todo-list/application"
	"github.com/sullyvannunes/todo-list/views"
)

func main() {
	view := views.New()

	app := application.New(view)

	log.Println("running on http://localhost:3030")
	http.ListenAndServe(":3030", app)
}
