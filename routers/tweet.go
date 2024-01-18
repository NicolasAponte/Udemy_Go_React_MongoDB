package routers

import (
	"context"
	"encoding/json"
	"time"

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

	if status {
		response.Message = "It wasn't possible save the tweet"
		return response
	}
	response.Status = 200
	response.Message = "Tweet saved successfully"
	return response
}
