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

	isOk, statusCode, msg, claim := validateAuthorization(ctx, request)
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
		case "tweet":
			return routers.SaveTweet(ctx, claim)
		case "avatar":
			return routers.UploadImage(ctx, "A", request, claim)
		case "banner":
			return routers.UploadImage(ctx, "B", request, claim)
		}
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		case "ping":
			return routers.Ping(ctx)
		case "getProfile":
			return routers.GetProfile(request)
		case "getTweets":
			return routers.GetTweets(request)
		case "getAvatar":
			return routers.GetImage(ctx, "A", request, claim)
		case "getBanner":
			return routers.GetImage(ctx, "B", request, claim)
		}
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		case "updateProfile":
			return routers.UpdateProfile(ctx, claim)
		}
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {
		case "deleteTweet":
			return routers.DeleteTweet(request, claim)
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

	token := request.Headers["Authorization"]

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
