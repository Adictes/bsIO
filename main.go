package main

import (
	"html/template"
	"log"
	"net/http"
)

var t *template.Template

// Context is struct for template
type Context struct {
	Title string
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", Greetings)

	t = template.Must(template.New("Game").ParseFiles("templates/index.html"))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Greetings says hello in browser's window
func Greetings(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index", Context{"bsIO"})
}
