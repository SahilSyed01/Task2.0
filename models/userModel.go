package models

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id"`
	First_name string             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  string             `json:"last_name" validate:"required,min=2,max=100"`
	Password   string             `json:"Password" validate:"required,min=6"`
	Email      string             `json:"email" validate:"email,required"`
	Phone      string             `json:"phone" validate:"required"`
	User_id    string             `json:"user_id"`
}

type UserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"Password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserID    string `json:"user_id"`
}

type SecretsManagerSecret struct {
	UserPoolID   string `json:"USER_POOL_ID"`
	ClientID     string `json:"CLIENT_ID"`
	ClientSecret string `json:"CLIENT_SECRET"`
	Username     string `json:"USERNAME"`
	Password     string `json:"PASSWORD"`
	Region       string `json:"REGION"`
}

type Secret struct {
	UserPoolID string `json:"user_pool_id"`
	Region     string `json:"region"`
}

type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}
