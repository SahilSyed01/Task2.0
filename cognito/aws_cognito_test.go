package cognito

import (
	"context"
	"testing"
	//"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/assert"
)

// MockCognitoIdentityProviderClient is a mock of the CognitoIdentityProviderClient interface
type MockCognitoIdentityProviderClient struct {
	InitiateAuthFn func(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error)
}

func (m *MockCognitoIdentityProviderClient) InitiateAuth(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return m.InitiateAuthFn(ctx, params, optFns...)
}

func TestGetJWTToken(t *testing.T) {
	mockClient := &MockCognitoIdentityProviderClient{
		InitiateAuthFn: func(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return &cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: &types.AuthenticationResultType{
					IdToken: aws.String("mockJWTToken"),
				},
			}, nil
		},
	}

	clientID := "testClientID"
	clientSecret := "testClientSecret"
	username := "testUser"
	password := "testPassword"
	userPoolID := "testUserPoolID"

	token, err := GetJWTToken(mockClient, userPoolID, clientID, clientSecret, username, password)
	assert.NoError(t, err)
	assert.Equal(t, "mockJWTToken", token)
}
