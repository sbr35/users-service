package main

import (
	"log"
	"net/http"
	"github.com/sbr35/wallets/services/users/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	log.Println("Listen on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
