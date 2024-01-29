package bd

import (
	"context"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertTweet(tweet models.SavedTweet) (string, bool, error) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("tweet")

	tweetTosave := bson.M{
		"user_id": tweet.UserID,
		"message": tweet.Message,
		"date":    tweet.Date,
	}

	result, err := col.InsertOne(ctx, tweetTosave)
	if err != nil {

		return "", false, err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.String(), true, nil
}

func SelectTweets(ID string, page int64) ([]*models.SelectedTweet, bool) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("tweet")

	var result []*models.SelectedTweet

	filter := bson.M{
		"user_id": ID,
	}

	options := options.Find()
	options.SetLimit(20)
	options.SetSort(bson.D{{Key: "date", Value: -1}})
	options.SetSkip((page - 1) * 20)

	cur, err := col.Find(ctx, filter, options)
	if err != nil {
		return result, false
	}

	for cur.Next(ctx) {
		var tweet models.SelectedTweet
		err := cur.Decode(&tweet)
		if err != nil {
			return result, false
		}
		result = append(result, &tweet)
	}

	return result, true
}

func DeleteTweet(ID string, UserID string) error {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("tweet")

	objID, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.M{
		"_id":     objID,
		"user_id": UserID,
	}

	_, err := col.DeleteOne(ctx, filter)
	return err
}
