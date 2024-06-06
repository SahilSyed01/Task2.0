package repository

import (
	"context"
	"time"
	"user-management-service/config"
	"user-management-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect() error {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	return err
}

func GetUserCollection() *mongo.Collection {
	return client.Database("userdb").Collection("users")
}

func CreateUser(user *models.User) error {
	_, err := GetUserCollection().InsertOne(context.Background(), user)
	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := GetUserCollection().FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	return &user, err
}
