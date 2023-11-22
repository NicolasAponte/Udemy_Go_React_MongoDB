package jwt

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

var (
	Email     string
	IDUsuario string
)

func ProcessToken(token string, JWTSign string) (*models.Claim, bool, string, error) {
	pass := []byte(JWTSign)
	var claims models.Claim

	splitToken := strings.Split(token, "bearer")
	if len(splitToken) != 2 {
		return &claims, false, "", errors.New("Invalid format token")
	}

	token = strings.TrimSpace(splitToken[1])
	tk, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return pass, nil
	})
	if err != nil {
		// Rutina contra BD
	}

	if !tk.Valid {
		return &claims, false, "", errors.New("Invalid Token")
	}

	return &claims, false, "", err
}
