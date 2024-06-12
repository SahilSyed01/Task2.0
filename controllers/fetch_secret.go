package controllers

import (
    "encoding/json"
    "fmt"
    // "log"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)


func getSecret(region, secretName string) (*SecretsManagerSecret, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess, &aws.Config{Region: aws.String(region)})
 
    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    }
 
    result, err := svc.GetSecretValue(input)
    if err != nil {
        return nil, err
    }
 
    if result.SecretString == nil {
        return nil, fmt.Errorf("secret string is nil")
    }
 
    secret := &SecretsManagerSecret{}
    err = json.Unmarshal([]byte(*result.SecretString), secret)
    if err != nil {
        return nil, err
    }
 
    return secret, nil
}
type SecretsManagerSecret struct {
    UserPoolID *string `json:"USER_POOL_ID"`
    Region     *string `json:"REGION"`
    // Add other fields from your secret here
}