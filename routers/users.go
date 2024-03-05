package routers

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func GetUserList(req events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	page := req.QueryStringParameters["page"]
	userType := req.QueryStringParameters["type"]
	search := req.QueryStringParameters["search"]

	UserID := claim.ID.Hex()

	if len(page) == 0 {
		page = "1"
	}

	pagTemp, err := strconv.Atoi(page)
	if err != nil {
		response.Message = "Page must be int" + err.Error()
		return response
	}

	users, status := bd.GetAllUsers(UserID, int64(pagTemp), search, userType)
	if !status {
		response.Message = "Error getting users"
		return response
	}

	respJson, err := json.Marshal(users)
	if err != nil {
		response.Status = 500
		response.Message = "Error formating users data " + err.Error()
		return response
	}

	response.Status = 200
	response.Message = string(respJson)

	return response
}
