package models

import (
	"REST-API/config"
	"REST-API/db"
	"REST-API/utils"
	"context"
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

func (u *User) Save(ctx context.Context) error {
	existingUser, err := GetUserByEmail(ctx, u.Email)
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

	result, err := db.DB.ExecContext(ctx, query, u.Email, string(hashedPassword), "user")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while creating user")
		}
		if ctx.Err() == context.Canceled {
			return errors.New("request was canceled while creating user")
		}
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

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password, role FROM users WHERE email = ?`

	row := db.DB.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("request timeout while fetching user")
		}
		return nil, err
	}

	return &user, nil
}

func (u *User) ValidateCredentials(ctx context.Context) error {
	existingUser, err := GetUserByEmail(ctx, u.Email)
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
func (u *User) SaveRefreshToken(ctx context.Context, token string) error {
	expiresAt := time.Now().Add(config.App.RefreshTokenExpiry)

	query := `INSERT INTO refresh_tokens(token, user_id, expires_at) VALUES (?, ?, ?)`

	_, err := db.DB.ExecContext(ctx, query, token, u.ID, expiresAt)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while saving refresh token")
		}
		return err
	}

	return nil
}

// checks if token exists, not expired
func ValidateRefreshToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password, u.role, rt.expires_at
		FROM users u
		INNER JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token = ?
	`

	var user User
	var expiresAt time.Time

	err := db.DB.QueryRowContext(ctx, query, token).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid refresh token")
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("request timeout while validating refresh token")
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
func DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`

	_, err := db.DB.ExecContext(ctx, query, token)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while deleting refresh token")
		}
		return err
	}

	return nil
}

// removes all refresh tokens for a user (for logout from all devices)
func (u *User) DeleteAllRefreshTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`

	_, err := db.DB.ExecContext(ctx, query, u.ID)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("request timeout while deleting refresh tokens")
		}
		return err
	}

	return nil
}

// invalidates old token and creates a new one
func (u *User) RotateRefreshToken(ctx context.Context, oldToken string) (string, error) {
	// Delete the old token
	err := DeleteRefreshToken(ctx, oldToken)
	if err != nil {
		return "", err
	}

	// Generate and save new token
	newToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	err = u.SaveRefreshToken(ctx, newToken)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
