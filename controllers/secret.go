// secret.go

package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var secretsClient *secretsmanager.Client

func init() {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    secretsClient = secretsmanager.NewFromConfig(cfg)
}

func getSecret(secretName string) (string, string, error) {
    input := &secretsmanager.GetSecretValueInput{
        SecretId: &secretName,
    }

    result, err := secretsClient.GetSecretValue(context.Background(), input)
    if err != nil {
        return "", "", err
    }

    var data map[string]string
    if err := json.Unmarshal([]byte(*result.SecretString), &data); err != nil {
        return "", "", err
    }

    region, ok := data["REGION"]
    if !ok {
        return "", "", fmt.Errorf("REGION not found in secret")
    }

    userPoolID, ok := data["USER_POOL_ID"]
    if !ok {
        return "", "", fmt.Errorf("USER_POOL_ID not found in secret")
    }

    return region, userPoolID, nil
}
