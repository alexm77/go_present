package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type user struct {
	Name     string    `json:"full_name"`
	Birthday time.Time `json:"date_of_birth"`
	Password string    `json:"password"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/super-endpoint", handle)
	router.HandleFunc("/super-endpoint/{id}", handleById)
	err := http.ListenAndServeTLS(":31337", "demo.crt", "demo.key", router)
	log.Fatal(err)
}

func handle(w http.ResponseWriter, r *http.Request) {
	println(r.Proto)
	user := user{"Super name", time.Now(), "you shall not pass!"}
	userBytes, _ := json.MarshalIndent(user, "", "  ")
	_, err := fmt.Fprintf(w, "I got your back\n%v\n", string(userBytes))
	if err != nil {
		log.Fatal(err)
	}
}
func handleById(w http.ResponseWriter, r *http.Request) {
	println(r.Proto)
	vars := mux.Vars(r)
	id := vars["id"]
	user := user{"Super name", time.Now(), "you shall not pass!"}
	userBytes, _ := json.MarshalIndent(user, "", "  ")
	_, err := fmt.Fprintf(w, "I got your back %s\n%v\n", id, string(userBytes))
	if err != nil {
		log.Fatal(err)
	}
}
