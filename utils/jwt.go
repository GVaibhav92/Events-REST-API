package utils

//Stateless authentication
import (
	"REST-API/config"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// creates a short-lived JWT access token
func GenerateToken(email string, userID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(config.App.AccessTokenExpiry).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return tokenString, nil
}

// validates a JWT access token and returns the user ID
func VerifyToken(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.App.JWTSecret), nil
	})

	if err != nil {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New("invalid token claims")
	}

	userID := int(claims["userId"].(float64))
	role := claims["role"].(string)
	return userID, role, nil
}

// creates a cryptographically secure random token
func GenerateRefreshToken() (string, error) {
	// create a 32-byte random token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes) // fills the byte slice with random data
	if err != nil {
		return "", errors.New("could not generate refresh token")
	}

	return hex.EncodeToString(bytes), nil
}
