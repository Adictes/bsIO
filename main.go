package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var (
	t           *template.Template
	store       *sessions.CookieStore
	curGames    map[string]string // Map: username to username, that now playing
	clean       map[string]bool
	fields      map[string]*Field // Game's fields assigned to users by their names
	shots       map[string]*Field // Players shots
	readyToPlay chan string       // User that ready to play
	turn        map[string]chan bool
	toSync      map[string]StrickenShips // map to synchronize StrickenShips
)

func init() {
	t = template.Must(template.New("Game").ParseFiles("templates/index.html", "templates/login.html"))
	store = sessions.NewCookieStore([]byte("very-secret-key"))
	curGames = make(map[string]string)
	clean = make(map[string]bool)
	fields = make(map[string]*Field)
	shots = make(map[string]*Field)
	readyToPlay = make(chan string, 1)
	turn = make(map[string]chan bool)
	toSync = make(map[string]StrickenShips)
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
	router.GET("/rff", RandomFieldFilling)
	router.GET("/clr", CleanAll)

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
	turn[session.Values["username"].(string)] = make(chan bool, 1)
	toSync[session.Values["username"].(string)] = StrickenShips{}
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
		fields[session.Values["username"].(string)].IndicateCell(msg[1], msg[3])
		ws.WriteJSON(fields[session.Values["username"].(string)].GetAvailableShips())
	}
}

// HitEnemyShips checks hitting on enemy's field, translates changes to both fields
// alters stricken ships each turn and checks whose turn is now
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
		if <-turn[session.Values["username"].(string)] == false {
			ws.WriteJSON(toSync[session.Values["username"].(string)])
			if as := fields[session.Values["username"].(string)].GetAvailableShips(); (as == Ships{4, 3, 2, 1}) {
				ws.WriteJSON(WinWrapper{false})
			}
			continue
		}

		ws.WriteJSON(toSync[session.Values["username"].(string)])
		ws.WriteJSON(TurnWrapper{true})

		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		enemy := GetEnemy(curGames, session.Values["username"].(string))
		if enemy == "" {
			// Значит не с кем поиграть, ждем
			turn[session.Values["username"].(string)] <- true
			continue
		}

		shots[session.Values["username"].(string)].IndicateCell(msg[1], msg[3])
		s := fields[enemy].GetStrickenShips(msg, session.Values["username"].(string))
		ws.WriteJSON(s)

		ChangeLetter(&s)
		toSync[enemy] = s

		if s.Hitted != "" {
			as := fields[enemy].GetAvailableShips()
			ws.WriteJSON(as)
			turn[session.Values["username"].(string)] <- true
			turn[enemy] <- false
			ws.WriteJSON(TurnWrapper{true})
			if (as == Ships{4, 3, 2, 1}) {
				clean[session.Values["username"].(string)] = true
				clean[enemy] = true
				ws.WriteJSON(WinWrapper{true})
			}
		} else {
			turn[enemy] <- true
			turn[session.Values["username"].(string)] <- false
			ws.WriteJSON(TurnWrapper{false})
		}

		//fields[enemy].print() // <-- for debug
	}
}

// StartTheGame initializes starting the game.
// First of all, it checks whether player can to play or not.
// He can't if he don't push the button 'I'm ready'.
// After that, it checks right positions of ships.
// If it's alright it adds player to chan "readyToPlay".
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
			ws.WriteJSON(CorrectnessWrapper{true})
		} else {
			ws.WriteJSON(CorrectnessWrapper{false})
			continue
		}

		log.Printf("User: %v want to play\n", session.Values["username"].(string))

		select {
		case un := <-readyToPlay:
			if enemy := HaveAvailableGame(curGames, un); enemy != "" {
				curGames[enemy] = un
				turn[un] <- false
				ws.WriteJSON(TurnWrapper{false})
			} else {
				curGames[un] = ""
				turn[un] <- true
				ws.WriteJSON(TurnWrapper{true})
			}
		}
		go func() {
			var enemy string
			for {
				enemy = GetEnemy(curGames, session.Values["username"].(string))
				if enemy != "" {
					ws.WriteJSON(NameWrapper{enemy})
					return
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

// RandomFieldFilling sets ships on the field randomly
func RandomFieldFilling(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := store.Get(r, "session") // Было session, err := (иначе не скомпилится)
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
		_, _, err := ws.ReadMessage() // На самом деле я не присваиваю пришедшее значение, так как в этом нет смысла
		if err != nil {
			log.Println(err)
			return
		}
		// Просто код ниже исполниться только тогда, когда
		// кнопку кто-то нажал, значит нужно сделать дело

		// Если что-то пришло, значит рандомно расставляем кораблики
		// даже если корабли уже расставлены, эта функция НИЖЕ будет работать
		// и переставлять кораблики
	}
}

// CleanAll cleans used vars and restore all to default
func CleanAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := store.Get(r, "session") // Было session, err := (иначе не скомпилится)
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
		// Тут пошла очистка
		for k := range clean {
			curGames[k] = ""
			fields[k] = &Field{}
			shots[k] = &Field{}
			turn[k] = make(chan bool, 1)
			toSync[k] = StrickenShips{}
			delete(clean, k)
		}
	}
}
