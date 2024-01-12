package routers

import (
	"context"
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

func UpdateProfile(ctx context.Context, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	var user models.User

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Error unmarshaling body " + err.Error()
		return response
	}

	status, err := bd.UpdateProfile(user, claim.ID.Hex())
	if err != nil {
		response.Message = "Error updating user " + err.Error()
		return response
	}

	if !status {
		response.Message = "Updating wasn't possible"
		return response
	}

	response.Status = 200
	response.Message = "Update Successful"
	return response

}
