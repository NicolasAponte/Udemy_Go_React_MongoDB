package bd

import (
	"context"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Registrate(user models.User) (string, bool, error) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("users")

	user.Password, _ = EncryptPass(user.Password)

	result, err := col.InsertOne(ctx, user)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}
