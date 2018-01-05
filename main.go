package main

import (
	"html/template"
	"log"
	"net/http"
)

var homeTmp *template.Template

// Context is struct for template
type Context struct {
	Title string
	Name  string
}

func main() {
	http.HandleFunc("/", Greetings)

	t := template.Must(template.New("Game").ParseFiles("template.html"))
	homeTmp = t.Lookup("template.html")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Greetings says hello in browser's window
func Greetings(w http.ResponseWriter, r *http.Request) {
	homeTmp.Execute(w, Context{"Doc", "Андрей"})
}
