package bd

import (
	"context"
	"fmt"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUsers(userID string, page int64, search string, userType string) ([]*models.User, bool) {
	ctx := context.TODO()

	db := MongoClient.Database(DatabaseName)
	col := db.Collection("users")

	var results []*models.User

	options := options.Find()
	options.SetLimit(20)
	options.SetSkip((page - 1) * 20)

	query := bson.M{
		"name": bson.M{
			"$regex": `(?i)` + search,
		},
	}

	cur, err := col.Find(ctx, query, options)
	if err != nil {
		return results, false
	}

	var include bool
	for cur.Next(ctx) {
		var s models.User
		err := cur.Decode(&s)
		if err != nil {
			fmt.Println("Decode = " + err.Error())
			return results, false
		}

		var r models.Relation
		r.UserID = userID
		r.UserRelationID = s.ID.Hex()

		include = false

		finded := GetRelation(r)
		if userType == "new" && !finded {
			include = true
		}
		if userType == "follow" && finded {
			include = true
		}

		if r.UserRelationID == userID {
			include = false
		}

		if include {
			s.Password = ""
			results = append(results, &s)
		}
	}

	err = cur.Err()
	if err != nil {
		fmt.Println("cur.Err = " + err.Error())
		return results, false
	}

	cur.Close(ctx)
	return results, true
}
