package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetAWSConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background())
}

func GetSecretsManagerClient(config aws.Config) *secretsmanager.Client {
	return secretsmanager.NewFromConfig(config)
}

func GetCognitoClient(config aws.Config) *cognitoidentityprovider.Client {
	return cognitoidentityprovider.NewFromConfig(config)
}
