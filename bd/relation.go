package bd

import (
	"context"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertRelation(relation models.Relation) (bool, error) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("relation")

	_, err := col.InsertOne(ctx, relation)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteRelation(relation models.Relation) (bool, error) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("relation")

	_, err := col.DeleteOne(ctx, relation)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetRelation(relation models.Relation) bool {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("relation")

	filter := bson.M{
		"user_id":          relation.UserID,
		"user_relation_id": relation.UserRelationID,
	}

	var result models.RelationValidate

	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return false
	}

	return true
}
