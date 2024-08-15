package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET string

func init() {
	JWT_SECRET = os.Getenv("JWT_SECRET")
}

func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println("invalid token")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
