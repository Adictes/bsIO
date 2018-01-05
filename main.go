package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Greetings)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Greetings says hello in browser's window
func Greetings(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, my dear friends!"))
}
