package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/naponte/Udemy_Go_React_MongoDB/awsgo"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/handlers"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	secretmanagergo "github.com/naponte/Udemy_Go_React_MongoDB/secretmanager.go"
)

func main() {
	/*
		La lambda en AWS se configuró bajo el nombre main, lo que indica que se va a llamar en este archivo main.go
	*/
	lambda.Start(EjecucionLambda)
}

// No se usa gin para gestionar la api, se usan funciones propias de AWS
func EjecucionLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InitAWS()

	if !ValidacionParams() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno",
			Headers: map[string]string{
				"Content-Type": "applicaton/json",
			},
		}
	}

	SecretModel, err := secretmanagergo.GetSecret(os.Getenv("SecretName"))

	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error in DB connection" + err.Error(),
			Headers: map[string]string{
				"Content-Type": "applicaton/json",
			},
		}
		return res, err
	}

	path := strings.Replace(request.PathParameters["twitterGo"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// Mongo Connection
	err = bd.InitConnection(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error en las varaibles de entorno",
			Headers: map[string]string{
				"Content-Type": "applicaton/json",
			},
		}
		return res, err
	}

	responseApi := handlers.Handlers(awsgo.Ctx, request)

	if responseApi.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: responseApi.Status,
			Body:       responseApi.Message,
			Headers: map[string]string{
				"Content-Type": "applicaton/json",
			},
		}
		return res, nil
	}

	return responseApi.CustomResp, nil
}

func ValidacionParams() bool {
	_, existSecretName := os.LookupEnv("SecretName")
	if !existSecretName {
		return false
	}

	_, existsBucketName := os.LookupEnv("BucketName")
	if !existsBucketName {
		return false
	}

	_, existUrlPrefix := os.LookupEnv("UrlPrefix")
	if !existUrlPrefix {
		return false
	}

	return true
}
