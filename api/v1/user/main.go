package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/gorilla/mux"
)

type Inode struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"file_name"`
	Path           string    `json:"filepath"`
	Type           string    `json:"file_type"`
	Mode           string    `json:"mode"`
	UID            uint64    `json:"uid"`
	GID            uint64    `json:"gid"`
	NumberOfBlocks uint64    `json:"num_of_blocks"`
	Size           uint64    `json:"size"`
	Timestamp      time.Time `json:"timestamp"`
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
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for _, inode := range inodes {
		if inode.ID == uint64(id) {
			// The fmt.Sprint() function in Go language formats
			//using the default formats for its operands and returns the resulting string.
			json.NewEncoder(w).Encode(inode)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	app := &application{}
	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/user", app.PostInode).Methods("POST")
	mux.HandleFunc("/api/v1/user", app.GetInode).Methods("GET")

	http.Handle("/", mux)

	port := 8080
	// use log.New() to create a logger for wr	iting information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// use stderr for writing error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.Printf("Starting server on %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	errorLog.Fatal(err)
}

//func init инициализируется до мейн функции
