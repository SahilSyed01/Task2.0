package config
 
import (
    "context"
    "testing"
 
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
    "github.com/stretchr/testify/assert"
)
 
// Mock implementations
type mockConfig struct{}
 
func (m *mockConfig) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
    return aws.Config{}, nil
}
 
type mockSecretsManager struct{}
 
func (m *mockSecretsManager) NewFromConfig(cfg aws.Config, optFns ...func(*secretsmanager.Options)) *secretsmanager.Client {
    return &secretsmanager.Client{}
}
 
type mockCognito struct{}
 
func (m *mockCognito) NewFromConfig(cfg aws.Config, optFns ...func(*cognitoidentityprovider.Options)) *cognitoidentityprovider.Client {
    return &cognitoidentityprovider.Client{}
}
 
// Test cases
func TestGetAWSConfig(t *testing.T) {
    originalConfigs := configs
    defer func() { configs = originalConfigs }()
 
    configs = (&mockConfig{}).LoadDefaultConfig
 
    cfg, err := GetAWSConfig()
    assert.NoError(t, err)
    assert.NotNil(t, cfg)
}
 
func TestGetSecretsManagerClient(t *testing.T) {
    originalSecrets := secrets
    defer func() { secrets = originalSecrets }()
 
    secrets = (&mockSecretsManager{}).NewFromConfig
 
    cfg := aws.Config{}
    client := GetSecretsManagerClient(cfg)
    assert.NotNil(t, client)
}
 
func TestGetCognitoClient(t *testing.T) {
    originalCognito := cognito
    defer func() { cognito = originalCognito }()
 
    cognito = (&mockCognito{}).NewFromConfig
 
    cfg := aws.Config{}
    client := GetCognitoClient(cfg)
    assert.NotNil(t, client)
}