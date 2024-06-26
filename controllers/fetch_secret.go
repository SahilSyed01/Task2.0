package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"go-chat-app/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretRetrievalError represents an error that occurred during secret retrieval.
type SecretRetrievalError struct {
	Message string
}

func (e SecretRetrievalError) Error() string {
	return fmt.Sprintf("Secret retrieval error: %s", e.Message)
}

// SecretsManagerClient is an interface for Secrets Manager client methods
type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var secretsManagerClient SecretsManagerClient

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	secretsManagerClient = secretsmanager.NewFromConfig(cfg)
}

func GetSecret(client SecretsManagerClient, secretName string) (*models.SecretsManagerSecret, error) {
	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, SecretRetrievalError{Message: err.Error()}
	}

	if result.SecretString == nil {
		return nil, SecretRetrievalError{Message: "secret string is nil"}
	}

	secret := &models.SecretsManagerSecret{}
	err = json.Unmarshal([]byte(*result.SecretString), secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
