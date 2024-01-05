package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func TokenGeneration(ctx context.Context, user models.User) (string, error) {
	jwtSign := ctx.Value(models.Key("jwtSign")).(string)
	tokenBytes := []byte(jwtSign)

	payload := jwt.MapClaims{
		"email":     user.Email,
		"name":      user.Name,
		"lastname":  user.Lastname,
		"birthdate": user.Biography,
		"biography": user.Biography,
		"location":  user.Location,
		"website":   user.WebSite,
		"_id":       user.ID.Hex(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(tokenBytes)
	if err != nil {
		return tokenStr, err
	}

	return tokenStr, nil
}
