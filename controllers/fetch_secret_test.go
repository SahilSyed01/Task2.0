// fetch_secret_test.go
package controllers

import (
    "context"
    "encoding/json"
    "errors"
    "testing"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
    "github.com/stretchr/testify/assert"
    "go-chat-app/models"
)

// MockSecretsManagerClient is a mock of the SecretsManagerClient interface
type MockSecretsManagerClient struct {
    GetSecretValueFunc func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// GetSecretValue implements SecretsManagerClient
func (m *MockSecretsManagerClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
    return m.GetSecretValueFunc(ctx, input, optFns...)
}

func TestGetSecret(t *testing.T) {
    originalSecretsManagerClient := secretsManagerClient
    defer func() { secretsManagerClient = originalSecretsManagerClient }()

    secretName := "testSecret"
    region := "us-west-2"

    validSecret := &models.SecretsManagerSecret{
        UserPoolID: "testPoolID",
        Region:     "us-west-2",
    }
    validSecretString, _ := json.Marshal(validSecret)

    secretsManagerClient = &MockSecretsManagerClient{
        GetSecretValueFunc: func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
            if *input.SecretId == secretName {
                return &secretsmanager.GetSecretValueOutput{
                    SecretString: aws.String(string(validSecretString)),
                }, nil
            }
            return nil, errors.New("some error")
        },
    }

    // Positive case: valid secret
    secret, err := GetSecret(region, secretName)
    assert.NoError(t, err)
    assert.NotNil(t, secret)
    assert.Equal(t, "testPoolID", secret.UserPoolID)
    assert.Equal(t, "us-west-2", secret.Region)

    // Negative case: secret retrieval error
    secretsManagerClient.(*MockSecretsManagerClient).GetSecretValueFunc = func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
        return nil, errors.New("some error")
    }

    secret, err = GetSecret(region, secretName)
    assert.Error(t, err)
    assert.Nil(t, secret)
    assert.IsType(t, SecretRetrievalError{}, err)
    assert.Equal(t, "Secret retrieval error: some error", err.Error())

    // Negative case: secret string is nil
    secretsManagerClient.(*MockSecretsManagerClient).GetSecretValueFunc = func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
        return &secretsmanager.GetSecretValueOutput{
            SecretString: nil,
        }, nil
    }

    secret, err = GetSecret(region, secretName)
    assert.Error(t, err)
    assert.Nil(t, secret)
    assert.IsType(t, SecretRetrievalError{}, err)
    assert.Equal(t, "Secret retrieval error: secret string is nil", err.Error())

    // Negative case: invalid JSON in secret string
    secretsManagerClient.(*MockSecretsManagerClient).GetSecretValueFunc = func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
        return &secretsmanager.GetSecretValueOutput{
            SecretString: aws.String("invalid json"),
        }, nil
    }

    secret, err = GetSecret(region, secretName)
    assert.Error(t, err)
    assert.Nil(t, secret)
    assert.Contains(t, err.Error(), "invalid character")
}
