package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var (
	t           *template.Template
	store       *sessions.CookieStore
	curGames    map[string]string // Map: username to username, that now playing
	fields      map[string]*Field // Game's fields assigned to users by their names
	shots       map[string]*Field
	readyToPlay chan string // User that ready to play
)

func init() {
	t = template.Must(template.New("Game").ParseFiles("templates/index.html", "templates/login.html"))
	store = sessions.NewCookieStore([]byte("very-secret-key"))
	curGames = make(map[string]string)
	fields = make(map[string]*Field)
	shots = make(map[string]*Field)
	readyToPlay = make(chan string, 1)
}

func main() {
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.GET("/", Index)
	router.GET("/login", LoginView)
	router.POST("/login", LoginSend)
	router.GET("/shs", SetHomeShips)
	router.GET("/hes", HitEnemyShips)
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
	shots[session.Values["username"].(string)] = &Field{}
	t.ExecuteTemplate(w, "index", session.Values["username"])
}

// LoginView displays login page
func LoginView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	t.ExecuteTemplate(w, "login", session.Values["username"])
}

// LoginSend sends username from form to session and authorized user
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

// SetHomeShips sets ships on the home field
// Uses websocket connection
func SetHomeShips(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		ws.WriteJSON(fields[session.Values["username"].(string)].GetAvailableShips())
	}
}

// StrickenShips is used as JSON wrapper for sending it to websocket
type StrickenShips struct {
	Ambient []string
	Hitted  string
}

// HitEnemyShips checks hit on enemy's field and send this data to websocket
func HitEnemyShips(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

		log.Printf("User: %v hit %v cell\n", session.Values["username"].(string), string(msg))

		enemy := FindEnemy(curGames, session.Values["username"].(string))
		if enemy == "" {
			// Значит не с кем поиграть, ждем
			continue
		}
		log.Printf("Вот такой противник - %v был найден игроку - %v\n", enemy, session.Values["username"].(string))

		shots[session.Values["username"].(string)].IndicateCell(msg[1], msg[3])
		if fields[enemy].Hit(msg[1], msg[3]) == true {
			flag, startRow, startCol, endRow, endCol := fields[enemy].isPadded(msg[1], msg[3], shots[session.Values["username"].(string)])
			startRow, startCol, endRow, endCol = startRow-1, startCol-1, endRow-1, endCol-1
			fmt.Println(startRow, startCol, endRow, endCol)
			if flag == true {
				s := StrickenShips{Hitted: string(msg)}
				if startCol == endCol {
					for i := startRow - 1; i <= endRow+1; i++ {
						s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", i, startCol-1))
						s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", i, startCol+1))
					}
					s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", startRow-1, startCol))
					s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", endRow+1, startCol))
				} else {
					for i := startCol - 1; i <= endCol+1; i++ {
						s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", startRow-1, i))
						s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", startRow+1, i))
					}
					s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", startRow, startCol-1))
					s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", startRow, endCol+1))
				}
				ws.WriteJSON(s)
				ws.WriteJSON(fields[enemy].GetAvailableShips())
			} else {
				ws.WriteJSON(StrickenShips{Hitted: string(msg)})
			}
		} else {
			ws.WriteJSON(StrickenShips{Ambient: []string{string(msg)}})
		}
		fields[enemy].print() // <-- for debug
	}
}

// StartTheGame initializes starting the game.
// First of all, it checks whether player can to play or not.
// He can't if he don't push the button 'I'm ready'.
// After that, it checks right positions of ships.
// If it's alright it adds player to readyToPlay chan.
// Then if it found player that also want to play,
// it adds him and they can to play.
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
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if fields[session.Values["username"].(string)].CheckPositionOfShips() == true {
			readyToPlay <- session.Values["username"].(string)
			ws.WriteJSON(true)
		} else {
			ws.WriteJSON(false)
			continue
		}

		log.Printf("User: %v want to play\n", session.Values["username"].(string))

		select {
		case un := <-readyToPlay:
			if enemy := HaveAvailableGame(curGames, un); enemy != "" {
				curGames[enemy] = un
			} else {
				curGames[un] = ""
			}
		}
	}
}
