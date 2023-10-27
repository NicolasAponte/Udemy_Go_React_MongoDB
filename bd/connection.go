package bd

import (
	"context"
	"fmt"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient  *mongo.Client
	DatabaseName string
)

func InitConnection(ctx context.Context) error {
	user := ctx.Value(models.Key("user")).(string)
	password := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)

	connStr := fmt.Sprintf("MongoDb+srv://%s:/%s@/%s?retryWrites=true&w=majority", user, password, host)

	var clientOptions = options.Client().ApplyURI(connStr)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Success MongoDB connection")
	MongoClient = client
	DatabaseName = ctx.Value(models.Key("database")).(string)

	return nil
}

func connectedDB() bool {
	err := MongoClient.Ping(context.TODO(), nil)
	return err == nil
}
