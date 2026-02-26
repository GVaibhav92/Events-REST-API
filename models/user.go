package models

import (
	"REST-API/config"
	"REST-API/db"
	"REST-API/utils"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	Role     string `json:"role"`
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

	query := `INSERT INTO users(email, password, role) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Email, string(hashedPassword), "user")
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(id)
	u.Role = "user"
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password, role FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
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
	u.Role = existingUser.Role
	return nil
}

// stores a refresh token in the database
func (u *User) SaveRefreshToken(token string) error {
	expiresAt := time.Now().Add(config.App.RefreshTokenExpiry)

	query := `INSERT INTO refresh_tokens(token, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(query, token, u.ID, expiresAt)

	return err
}

// checks if token exists, not expired
func ValidateRefreshToken(token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password,u.role,rt.expires_at
		FROM users u
		INNER JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token = ?
	`

	var user User
	var expiresAt time.Time

	err := db.DB.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return nil, errors.New("refresh token expired")
	}

	return &user, nil
}

// removes a specific refresh token (for logout)
func DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := db.DB.Exec(query, token)
	return err
}

// removes all refresh tokens for a user (for logout from all devices)
func (u *User) DeleteAllRefreshTokens() error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := db.DB.Exec(query, u.ID)
	return err
}

// invalidates old token and creates a new one
func (u *User) RotateRefreshToken(oldToken string) (string, error) {
	// Delete the old token
	err := DeleteRefreshToken(oldToken)
	if err != nil {
		return "", err
	}

	// Generate and save new token
	newToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	err = u.SaveRefreshToken(newToken)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
