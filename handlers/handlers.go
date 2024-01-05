package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/naponte/Udemy_Go_React_MongoDB/jwt"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"github.com/naponte/Udemy_Go_React_MongoDB/routers"
	"github.com/naponte/Udemy_Go_React_MongoDB/utils"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.Response {
	fmt.Println("Will be processed " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))
	var response models.Response
	response.Status = 400

	isOk, statusCode, msg, _ := validateAuthorization(ctx, request) //claim
	if !isOk {
		response.Status = statusCode
		response.Message = msg
		return response
	}

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "registration":
			return routers.Registration(ctx)
		case "login":
			return routers.Login(ctx)
		}
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		case "ping":
			return routers.Ping(ctx)
		}
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {

		}
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {

		}
	}

	response.Message = "Invalid Method"
	return response
}

func validateAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	noAutorization := []string{"registration", "login", "getAvatar", "getBanner", "ping"}
	path := ctx.Value(models.Key("path")).(string)
	if utils.SliceContains(path, noAutorization) {
		return true, http.StatusOK, "", models.Claim{}
	}

	token := request.Headers["Autorization"]

	if len(token) == 0 {
		return false, http.StatusUnauthorized, "Token is required", models.Claim{}
	}

	claim, isOk, msg, err := jwt.ProcessToken(token, ctx.Value(models.Key("jwtSign")).(string))
	if !isOk {
		if err != nil {
			fmt.Println("Error in token content " + err.Error())
			return false, http.StatusUnauthorized, err.Error(), models.Claim{}
		} else {
			fmt.Println("Error in token " + msg)
			return false, http.StatusUnauthorized, msg, models.Claim{}
		}
	}

	fmt.Println("Token OK")
	return true, http.StatusOK, msg, *claim
}
