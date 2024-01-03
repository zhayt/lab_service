package main

import "database/sql"

var db *sql.DB

func dbCreateUser(user User) (uint64, error) {
	res, err := db.Exec("INSERT INTO web_user (name, surname, email, password, type) VALUES (?, ?, ?, ?, ?)", user.Name, user.Surname, user.Email, user.Password, user.Type)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func dbUpdateUserPassword(user UserChangePasswordDTO) error {
	_, err := db.Exec("UPDATE web_user SET password=?, updated_at = NOW() WHERE email=?", user.NewPassword, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func dbGetUser(id uint64) (*User, error) {
	user := User{}
	if err := db.QueryRow("SELECT * FROM web_user WHERE id=? AND deleted_at IS NULL", id).Scan(&user.Name, &user.Surname, &user.Email, &user.Password, &user.Type, &user.DeletedAt, &user.UpdatedAt, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func dbGetUserByEmail(email string) (User, error) {
	var user = &User{}
	if err := db.QueryRow("SELECT email, password FROM web_user WHERE email=?", email).Scan(&user.Email, &user.Password); err != nil {
		return User{}, err
	}
	return *user, nil
}

func dbGetUsers() (*[]User, error) {
	row, err := db.Query("SELECT * FROM web_user")
	if err != nil {
		return nil, err
	}
	users := []User{}
	for row.Next() {
		var user User
		if err := row.Scan(&user.Name, &user.Surname, &user.Email, &user.Password, &user.Type, &user.DeletedAt, &user.UpdatedAt, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func dbDeleteUser(id uint64) error {
	if _, err := db.Exec("UPDATE web_user SET deleted_at=NOW() WHERE id=?", id); err != nil {
		return err
	}

	return nil
}
