package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/gorilla/mux"
)

type Inode struct {
	ID               int       `json:"id"`
	Name             string    `json:"file_name"`
	Path             string    `json:"filepath"`
	Type             string    `json:"file_type"`
	Mode             string    `json:"mode"`
	Uid              int       `json:"uid"`
	Gid              int       `json:"gid"`
	Number_of_Blocks int       `json:"num_of_blocks"`
	Size             int       `json:"size"`
	Timestamp        time.Time `json:"timestamp"`
}

var inodes []Inode

func (app *application) PostInode(w http.ResponseWriter, r *http.Request) {
	var newStruct Inode
	err := json.NewDecoder(r.Body).Decode(&newStruct)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	newStruct.Timestamp = time.Now()
	inodes = append(inodes, newStruct)

	w.WriteHeader(http.StatusCreated)
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
}

func (app *application) GetInode(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	for _, inode := range inodes {
		if fmt.Sprint(inode.ID) == id {
			json.NewEncoder(w).Encode(inode)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	app := &application{}
	mux := mux.NewRouter()
	mux.HandleFunc("/save", app.PostInode).Methods("POST")
	mux.HandleFunc("/parse", app.GetInode).Methods("GET")

	http.Handle("/", mux)

	port := 8080
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	// use log.New() to create a logger for wr	iting information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// use stderr for writing error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.Printf("Starting server on %d", 8080)
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
	errorLog.Fatal(err)
}
