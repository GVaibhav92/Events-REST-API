package models

import (
	"REST-API/db"
	"REST-API/utils"
	"database/sql"
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

func (u *User) Save() error {
	existingUser, err := GetUserByEmail(u.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users(email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Email, string(hashedPassword))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(id)
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *User) ValidateCredentials() error {
	existingUser, err := GetUserByEmail(u.Email)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("invalid credentials")
	}

	passwordMatch := utils.CheckPasswordHash(u.Password, existingUser.Password)
	if !passwordMatch {
		return errors.New("invalid credentials")
	}

	u.ID = existingUser.ID
	u.Password = existingUser.Password
	return nil
}
