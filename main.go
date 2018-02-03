package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var (
	store  *sessions.CookieStore
	t      *template.Template
	fields map[string]*Field
	bot    Field
)

func init() {
	store = sessions.NewCookieStore([]byte("very-secret-key"))
	t = template.Must(template.New("Game").ParseFiles("templates/index.html", "templates/login.html"))
	fields = make(map[string]*Field)
}

func main() {
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.GET("/", Index)
	router.GET("/login", LoginView)
	router.POST("/login", LoginSend)
	router.GET("/ps", PressCell)
	router.GET("/stg", StartTheGame)

	err := http.ListenAndServe(":8080", context.ClearHandler(router))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Index is general page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	if auth, ok := session.Values["logged"].(bool); !auth || !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	fields[session.Values["username"].(string)] = &Field{}
	t.ExecuteTemplate(w, "index", session.Values["username"])
}

// LoginView views login page
func LoginView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	t.ExecuteTemplate(w, "login", session.Values["username"])
}

// LoginSend sends data's form to session and authorized user
func LoginSend(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	if err := r.ParseForm(); err != nil {
		log.Fatal("Form parsing: ", err)
	}

	session.Values["username"] = r.FormValue("username")
	session.Values["logged"] = true
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// PressCell is action, when user press the cell for putting ship in this one.
// Uses websocket connection
func PressCell(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

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
		fields[session.Values["username"].(string)].IndicateCell(msg[0], msg[2])

		//ws.WriteJSON(fields[session.Values["username"].(string)].GetNotAccessibleCells())
		ws.WriteJSON(fields[session.Values["username"].(string)].GetAvailableShips())
	}
}

// StartTheGame checks map validation and if it's alright searches available user
// and starts game with him
func StartTheGame(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

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
		log.Println(string(msg) + " by " + session.Values["username"].(string))

		// Here function that checks validation

		// Here selection of players
	}
}
