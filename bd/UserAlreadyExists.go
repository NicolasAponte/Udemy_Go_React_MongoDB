package bd

import (
	"context"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
)

func UserAlreadyExists(email string) (models.User, bool, string) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("users")

	condition := bson.M{"email": email}

	var result models.User

	err := col.FindOne(ctx, condition).Decode(&result)
	ID := result.ID.Hex()
	if err != nil {
		return result, false, ID
	}

	return result, true, ID
}
