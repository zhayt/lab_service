package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Name      string     `json:"name,omitempty"`
	Surname   string     `json:"surname,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	Type      string     `json:"type,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
}

type UserChangePasswordDTO struct {
	Email       string `json:"email,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

type UserCreateResponse struct {
	UserID          uint64 `json:"user_id"`
	ResponseMessage string `json:"response_message"`
}

var salt = "wenvc32nocoj1313dc#@D#vc"

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Printf("Parse request body error: %s", err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// TODO: validate data
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		logger.Printf("Hash password error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	user.Password = string(passwordHash)
	userID, err := dbCreateUser(user)
	if err != nil {
		logger.Printf("Create user error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

	response := UserCreateResponse{UserID: userID, ResponseMessage: "User successfully created"}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Printf("Send response error: %s", err.Error())
		return
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	userID, err := strconv.ParseUint(r.URL.Query().Get("id"), 64, 0)
	if err != nil {
		logger.Printf("Parse userID from url param error: %s", err.Error())
		http.Error(w, http.StatusText(404), 404)
		return
	}
	user, err := dbGetUser(userID)
	if err != nil {
		logger.Printf("Get user error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.Printf("Send user data response error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	users, err := dbGetUsers()
	if err != nil {
		logger.Printf("Get users error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		logger.Printf("Send users as response error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	var user UserChangePasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Printf("Parse request body error: %s", err.Error())
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	// TODO: validate data
	savedUser, err := dbGetUserByEmail(user.Email)
	if err != nil {
		logger.Printf("Getting user error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.OldPassword+salt)); err != nil {
		logger.Printf("Compare user password error: %s", err.Error())
		http.Error(w, http.StatusText(400), 400)
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword+salt), bcrypt.DefaultCost)
	if err != nil {
		logger.Printf("Hash password error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	user.NewPassword = string(hasedPassword)
	if err := dbUpdateUserPassword(user); err != nil {
		logger.Printf("Update user password error: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

// TODO: implement it
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
