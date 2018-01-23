package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var t *template.Template

func main() {
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.GET("/", Index)
	router.GET("/ws", PressCell)

	t = template.Must(template.New("Game").ParseFiles("templates/index.html"))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// PressCell is action, when user press the cell for putting ship in this one.
// Uses websocket connection
func PressCell(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(msg))
		IndicateCell(msg[0], msg[2])

		//ws.WriteJSON(GetNotAccessibleCells())
		ws.WriteJSON(GetAvailableShips())
	}
}

// Index is general page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fieldInit()
	t.ExecuteTemplate(w, "index", nil)
}
