package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/gorilla/mux"
)

type Inode struct {
	ID             uint64    `json:"id"`
	FileName       string    `json:"file_name"`
	Path           string    `json:"filepath"`
	Type           string    `json:"file_type"`
	Mode           string    `json:"mode"`
	UID            uint64    `json:"uid"`
	GID            uint64    `json:"gid"`
	NumberOfBlocks uint64    `json:"num_of_blocks"`
	Size           uint64    `json:"size"`
	Timestamp      time.Time `json:"timestamp"`
}

var inodes = make([]Inode, 0, 1000)
var InodesMap = make(map[uint64]Inode)

func (app *application) SaveInode(w http.ResponseWriter, r *http.Request) {
	var NewInode Inode
	err := json.NewDecoder(r.Body).Decode(&NewInode)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	NewInode.Timestamp = time.Now()

	var mu sync.Mutex
	mu.Lock()
	inodes = append(inodes, NewInode)
	InodesMap[NewInode.ID] = NewInode
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
}

func (app *application) ParseInode(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	for _, inode := range inodes {
		if inode.ID == uint64(id) {
			json.NewEncoder(w).Encode(inode)
			return
		}
	}
	inode, exists := InodesMap[uint64(id)]
	if exists {
		json.NewEncoder(w).Encode(inode)
		return
	}
	http.NotFound(w, r)
}

func main() {
	app := &application{}
	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/user", app.SaveInode).Methods("POST")
	mux.HandleFunc("/api/v1/user", app.ParseInode).Methods("GET")

	http.Handle("/", mux)

	port := 8080
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.Printf("Starting server on %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		errorLog.Printf("Error is discovered %s", err)
		return
	}
}

//func init инициализируется до мейн функции
