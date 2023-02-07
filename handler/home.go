package handler

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
)

func (h Handler) Home(w http.ResponseWriter, r *http.Request) {
	var err error
	t := template.New("Welcome Home")
	t = template.Must(t.ParseFiles("templates/home.html"))
	if err != nil {
		log.Fatalf("%v", err)
	}

	t.ExecuteTemplate(w, "home.html" ,nil)
}
