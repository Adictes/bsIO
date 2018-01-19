package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var t *template.Template

// Context is struct for template
type Context struct {
	Title string
}

func main() {
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.GET("/", Index)
	router.POST("/", PressCell)

	t = template.Must(template.New("Game").ParseFiles("templates/index.html"))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// PressCell is action, when user press the dot for putting ship in this cell
func PressCell(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	log.Println(r.Form["pos"])
	t.ExecuteTemplate(w, "index", Context{"bsIO"})
}

// Index is general page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t.ExecuteTemplate(w, "index", Context{"bsIO"})
}
