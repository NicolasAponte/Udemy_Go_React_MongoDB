package bd

import (
	"context"
	"fmt"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SearchProfile(ID string) (models.User, error) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("users")

	var profile models.User

	objID, _ := primitive.ObjectIDFromHex(ID)

	condition := bson.M{
		"_id": objID,
	}

	err := col.FindOne(ctx, condition).Decode(&profile)
	if err != nil {
		return profile, err
	}

	profile.Password = ""

	return profile, nil
}

func UpdateProfile(user models.User, ID string) (bool, error) {
	ctx := context.TODO()
	db := MongoClient.Database(DatabaseName)
	col := db.Collection("users")

	profile := make(map[string]interface{})
	if len(user.Name) > 0 {
		profile["name"] = user.Name
	}
	if len(user.Lastname) > 0 {
		profile["lastname"] = user.Name
	}
	if len(user.Name) > 0 {
		profile["name"] = user.Name
	}

	profile["birthdate"] = user.BirthDate

	if len(user.Avatar) > 0 {
		profile["avatar"] = user.Avatar
	}

	if len(user.Banner) > 0 {
		profile["banner"] = user.Banner
	}
	if len(user.Biography) > 0 {
		profile["biography"] = user.Biography
	}
	if len(user.Location) > 0 {
		profile["location"] = user.Location
	}
	if len(user.WebSite) > 0 {
		profile["website"] = user.WebSite
	}

	fmt.Println(profile)
	updateStr := bson.M{
		"$set": profile,
	}

	objID, _ := primitive.ObjectIDFromHex(ID)
	fmt.Println(objID)
	filter := bson.M{
		"_id": bson.M{
			"$eq": objID,
		},
	}

	result, err := col.UpdateOne(ctx, filter, updateStr)
	fmt.Println(result)
	if err != nil {
		return false, err
	}

	return true, nil
}
