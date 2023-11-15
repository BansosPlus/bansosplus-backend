package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"

	"github.com/BansosPlus/bansosplus-backend.git/model"
)

func GenerateToken(user *model.User) (string, error) {
	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign the token with a secret key
	secretKey := []byte("bansosplus-secret-key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte("bansosplus-secret-key")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Validate the token
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}