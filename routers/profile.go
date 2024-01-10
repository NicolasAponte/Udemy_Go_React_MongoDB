package routers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func GetProfile(req events.APIGatewayProxyRequest) models.Response {
	var response models.Response
	response.Status = 400

	fmt.Println("Into profile")

	ID := req.QueryStringParameters["id"]
	if len(ID) < 1 {
		response.Message = "You have to send ID"
		return response
	}

	profile, err := bd.SearchProfile(ID)
	if err != nil {
		response.Message = "Error searching profile " + err.Error()
		return response
	}

	respJson, err := json.Marshal(profile)
	if err != nil {
		response.Message = "Error formatting profile to JSON " + err.Error()
		response.Status = 500
		return response
	}

	response.Status = 200
	response.Message = string(respJson)
	return response
}
