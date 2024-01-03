package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var server http.Server
var logger *log.Logger

func main() {
	fmt.Println("Starting app")
	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/user", createUserHandler).Methods("POST")
	mux.HandleFunc("/api/v1/user/{id}", getUserHandler).Methods("GET")
	mux.HandleFunc("/api/v1/user", getUsersHandler).Methods("GET")
	mux.HandleFunc("/api/v1/user", updateUserPasswordHandler).Methods("PATCH")
	mux.HandleFunc("/api/v1/user", deleteUserHandler).Methods("DELETE")
	server.Handler = mux
	if err := server.ListenAndServe(); err != nil {
		logger.Printf("Run server error: %s", err.Error())
		log.Fatal(err)
	}
}
