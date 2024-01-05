package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/jwt"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func Login(ctx context.Context) models.Response {
	var user models.User
	var response models.Response
	response.Status = 400

	body := ctx.Value(models.Key("body")).(string)

	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Unable to read body: " + err.Error()
		return response
	}

	if len(user.Email) == 0 {
		response.Message = "Email is required"
		return response
	}

	_, exists := bd.Login(user.Email, user.Password)
	if !exists {
		response.Message = "Invalid Email or Password"
		return response
	}

	jwtKey, err := jwt.TokenGeneration(ctx, user)
	if err != nil {
		response.Message = "Error in token generation: " + err.Error()
		return response
	}

	responseLogin := models.LoginResponse{
		Token: jwtKey,
	}

	token, err := json.Marshal(responseLogin)
	if err != nil {
		response.Message = "Error formatting token: " + err.Error()
		return response
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(time.Hour * 24),
	}

	cookieStr := cookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow_Origin": "*",
			"Set-Cookie":                  cookieStr,
		},
	}

	response.Status = 200
	response.Message = string(token)
	response.CustomResp = res

	return response
}
