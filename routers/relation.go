package routers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func CreateRelation(ctx context.Context, req events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	ID := req.QueryStringParameters["id"]
	if len(ID) < 1 {
		response.Message = "ID is requered"
		return response
	}

	var relation models.Relation
	relation.UserID = claim.ID.Hex()
	relation.UserRelationID = ID

	status, err := bd.InsertRelation(relation)
	if err != nil {
		response.Message = "Error saving relation in db " + err.Error()
		return response
	}

	if !status {
		response.Message = "It was not possible save the relation"
		return response
	}

	response.Status = 200
	response.Message = "Relation saved successfully"

	return response
}

func DeleteRelation(req events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	ID := req.QueryStringParameters["id"]
	if len(ID) < 1 {
		response.Message = "ID is requered"
		return response
	}

	var relation models.Relation
	relation.UserID = claim.ID.Hex()
	relation.UserRelationID = ID

	status, err := bd.DeleteRelation(relation)
	if err != nil {
		response.Message = "Error dropping relation in db " + err.Error()
		return response
	}

	if !status {
		response.Message = "It was not possible delete the relation"
		return response
	}

	response.Status = 200
	response.Message = "Relation deleted successfully"

	return response
}

func GetRelation(req events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	ID := req.QueryStringParameters["id"]
	if len(ID) < 1 {
		response.Message = "ID is requered"
		return response
	}

	var relation models.Relation
	relation.UserID = claim.ID.Hex()
	relation.UserRelationID = ID

	var relationVal models.RelationValidate
	relationVal.Status = false

	existRelation := bd.GetRelation(relation)
	if existRelation {
		relationVal.Status = true
	}

	respJson, err := json.Marshal(existRelation)
	if err != nil {
		response.Status = 500
		response.Message = "Error formating relation data " + err.Error()
		return response
	}

	response.Status = 200
	response.Message = string(respJson)

	return response
}
