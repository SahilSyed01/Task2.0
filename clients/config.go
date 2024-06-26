package config
 
import (
    "context"
 
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)
 
var(
    configs=config.LoadDefaultConfig
    secrets=secretsmanager.NewFromConfig
    cognito=cognitoidentityprovider.NewFromConfig
)
 
func GetAWSConfig() (aws.Config, error) {
    return configs(context.Background())
}
 
func GetSecretsManagerClient(config aws.Config) *secretsmanager.Client {
    return secrets(config)
}
 
func GetCognitoClient(config aws.Config) *cognitoidentityprovider.Client {
    return cognito(config)
}