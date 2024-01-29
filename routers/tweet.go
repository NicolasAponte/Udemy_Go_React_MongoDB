package routers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func SaveTweet(ctx context.Context, claim models.Claim) models.Response {
	var message models.Tweet
	var response models.Response
	response.Status = 400

	userID := claim.ID.Hex()

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &message)
	if err != nil {
		response.Message = "Error unmarshaling message"
		return response
	}

	tweet := models.SavedTweet{
		UserID:  userID,
		Message: message.Message,
		Date:    time.Now(),
	}

	_, status, err := bd.InsertTweet(tweet)
	if err != nil {
		response.Message = "Error saving tweet " + err.Error()
		return response
	}

	if !status {
		response.Message = "It wasn't possible save the tweet"
		return response
	}
	response.Status = 200
	response.Message = "Tweet saved successfully"
	return response
}

func GetTweets(request events.APIGatewayProxyRequest) models.Response {
	var response models.Response
	response.Status = 400

	ID := request.QueryStringParameters["id"]
	page := request.QueryStringParameters["page"]

	if len(ID) < 1 {
		response.Message = "ID parameter is required"
		return response
	}

	if len(page) < 1 {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		response.Message = "Page parameter must be numeric"
		return response
	}

	tweets, status := bd.SelectTweets(ID, int64(pageInt))
	if !status {
		response.Message = "Error getting page"
	}

	jsonResponse, err := json.Marshal(tweets)
	if err != nil {
		response.Message = "Error marshalling tweet " + err.Error()
		response.Status = 500
		return response
	}

	response.Status = 200
	response.Message = string(jsonResponse)
	return response
}
