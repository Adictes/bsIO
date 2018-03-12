package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	fields = make(map[string]*Field)
	shots = make(map[string]*Field)
	readyToPlay = make(chan string, 1)
	turn = make(map[string]chan bool)
	toSync = make(map[string]StrickenShips)

	store.Options = &sessions.Options{
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
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

// upgrader uses to establish websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// SetHomeShips sets ships on the home field
func SetHomeShips(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrading:", err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message:", err)
			continue
		}

		fields[session.Values["username"].(string)].IndicateCell(msg[1], msg[3])
		ws.WriteJSON(fields[session.Values["username"].(string)].GetAvailableShips())
	}
}

// HitEnemyShips checks hitting on enemy's field,
// translates changes to both fields, alters stricken ships each turn,
// checks whose turn is now and tracking for win/lose
func HitEnemyShips(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrading:", err)
		return
	}
	defer ws.Close()

	for {
		if <-turn[session.Values["username"].(string)] == false {
			ws.WriteJSON(toSync[session.Values["username"].(string)])
			ws.WriteJSON(TurnWrapper{false})
			if as := fields[session.Values["username"].(string)].GetAvailableShips(); (as == Ships{4, 3, 2, 1}) {
				ws.WriteJSON(WinWrapper{false})
			}
			continue
		}

		ws.WriteJSON(toSync[session.Values["username"].(string)])
		ws.WriteJSON(TurnWrapper{true})

		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message:", err)
			continue
		}

		enemy := GetEnemy(curGames, session.Values["username"].(string))
		if enemy == "" {
			// Значит не с кем поиграть, ждем
			turn[session.Values["username"].(string)] <- true
			continue
		}

		// Записали попадание в shots, после получили структуру
		// сбитых и прилежащих к сбитым ячеек, отправили ее для рендера
		// у себя, изменили ее для последующего рендера у 2-ого игрока
		// и записали эти изменения на его имя
		shots[session.Values["username"].(string)].IndicateCell(msg[1], msg[3])
		strickenShips := fields[enemy].GetStrickenShips(msg, session.Values["username"].(string))
		ws.WriteJSON(strickenShips)
		ChangeLetter(&strickenShips)
		toSync[enemy] = strickenShips

		if strickenShips.Hitted != "" {
			as := fields[enemy].GetAvailableShips()
			ws.WriteJSON(as)
			turn[enemy] <- false
			if (as == Ships{4, 3, 2, 1}) {
				ws.WriteJSON(WinWrapper{true})
				continue
			}
			turn[session.Values["username"].(string)] <- true
			ws.WriteJSON(TurnWrapper{true})
		} else {
			turn[enemy] <- true
			turn[session.Values["username"].(string)] <- false
			ws.WriteJSON(TurnWrapper{false})
		}
	}
}

// StartTheGame initializes starting the game.
// First of all, it checks right positions of ships.
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
		log.Println("WebSocket upgrading:", err)
		return
	}
	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message:", err)
			continue
		}

		if fields[session.Values["username"].(string)].CheckPositionOfShips() == true {
			readyToPlay <- session.Values["username"].(string)
			ws.WriteJSON(CorrectnessWrapper{true})
		} else {
			ws.WriteJSON(CorrectnessWrapper{false})
			continue
		}

		select {
		case un := <-readyToPlay:
			if enemy := HaveAvailableGame(curGames, un); enemy != "" {
				curGames[enemy] = un
				turn[un] <- false
				ws.WriteJSON(TurnWrapper{false})
				log.Printf("User: %v is starting to play with %v\n", un, enemy)
			} else {
				curGames[un] = ""
				turn[un] <- true
				ws.WriteJSON(TurnWrapper{true})
				log.Printf("User: %v is waiting to play with someone\n", un)
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
	rand.Seed(time.Now().UnixNano())
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatal("Session: ", err)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrading:", err)
		return
	}
	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message:", err)
			continue
		}

		ws.WriteJSON(ClearWrapper{true})

		var tempField Field
		fields[session.Values["username"].(string)] = &Field{}
		for len := 4; len >= 1; len-- { //цикл по длине кораблей
			for k := len; k <= 4; k++ { //цикл по кол-ву кораблей
				flag := true
				var row1, col1, row2, col2 uint8
				for flag == true {
					flag = false
					row1 = uint8(rand.Int()%(10-len+1)) + 1
					col1 = uint8(rand.Int()%(10-len+1)) + 1
					if orientation := uint8(rand.Int() % 2); orientation == 1 {
						col2 = col1 + uint8(len-1)
						row2 = row1
					} else {
						row2 = row1 + uint8(len-1)
						col2 = col1
					}
					//проверяем что наш корабль не пересекается с уже заданными
					for i := col1 - 1; i <= col2+1; i++ {
						for j := row1 - 1; j <= row2+1; j++ {
							if tempField[j][i] == true {
								flag = true
							}
						}
					}
				}
				//если мы тут, значит корабль можно поставить
				//заносим его в наш массив(поле):
				for i := col1; i <= col2; i++ {
					for j := row1; j <= row2; j++ {
						tempField[j][i] = true
					}
				}
			}
		}
		for i := 1; i <= 10; i++ {
			for j := 1; j <= 10; j++ {
				if tempField[i][j] == true {
					fields[session.Values["username"].(string)].IndicateCell(byte('0'+i-1), byte('0'+j-1))
					ws.WriteJSON(CellWrapper{"h" + strconv.Itoa(i-1) + "-" + strconv.Itoa(j-1)})
				}
			}
		}

		ws.WriteJSON(fields[session.Values["username"].(string)].GetAvailableShips())
	}
}

// CleanAll cleans used vars and restore all to default
func CleanAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Println("Session: ", err)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrading:", err)
		return
	}
	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message:", err)
			continue
		}

		delete(curGames, session.Values["username"].(string))
		fields[session.Values["username"].(string)] = &Field{}
		shots[session.Values["username"].(string)] = &Field{}
		toSync[session.Values["username"].(string)] = StrickenShips{}
	}
}
