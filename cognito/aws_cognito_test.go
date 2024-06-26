 
package cognito
 
import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "testing"
 
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
    if m.InitiateAuthFn != nil {
        return m.InitiateAuthFn(ctx, params, optFns...)
    }
    return &cognitoidentityprovider.InitiateAuthOutput{
        AuthenticationResult: &types.AuthenticationResultType{
            AccessToken: aws.String("access-token"),
            ExpiresIn:   int32(20),
        },
    }, nil
}
 
func TestGetJWTToken(t *testing.T) {
    t.Run("successful auth", func(t *testing.T) {
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
 
        computeSecretHash = func(clientSecret, clientID, username string) string {
            return "mockSecretHash"
        }
 
        token, err := GetJWTToken(mockClient, userPoolID, clientID, clientSecret, username, password)
        assert.NoError(t, err)
        assert.Equal(t, "mockJWTToken", token)
    })
 
    t.Run("auth error", func(t *testing.T) {
        mockClient := &MockCognitoIdentityProviderClient{
            InitiateAuthFn: func(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error) {
                return nil, fmt.Errorf("mock auth error")
            },
        }
 
        clientID := "testClientID"
        clientSecret := "testClientSecret"
        username := "testUser"
        password := "testPassword"
        userPoolID := "testUserPoolID"
 
        computeSecretHash = func(clientSecret, clientID, username string) string {
            return "mockSecretHash"
        }
 
        token, err := GetJWTToken(mockClient, userPoolID, clientID, clientSecret, username, password)
        assert.Error(t, err)
        assert.Equal(t, "", token)
        assert.Contains(t, err.Error(), "failed to initiate auth")
    })
 
    t.Run("nil auth result", func(t *testing.T) {
        mockClient := &MockCognitoIdentityProviderClient{
            InitiateAuthFn: func(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.InitiateAuthOutput, error) {
                return &cognitoidentityprovider.InitiateAuthOutput{
                    AuthenticationResult: nil,
                }, nil
            },
        }
 
        clientID := "testClientID"
        clientSecret := "testClientSecret"
        username := "testUser"
        password := "testPassword"
        userPoolID := "testUserPoolID"
 
        computeSecretHash = func(clientSecret, clientID, username string) string {
            return "mockSecretHash"
        }
 
        token, err := GetJWTToken(mockClient, userPoolID, clientID, clientSecret, username, password)
        assert.Error(t, err)
        assert.Equal(t, "", token)
        assert.Contains(t, err.Error(), "authentication result is nil")
    })
}
 
func TestComputeSecretHash(t *testing.T) {
    clientSecret := "testClientSecret"
    clientID := "testClientID"
    username := "testUser"
 
    expectedHash := func(clientSecret, clientID, username string) string {
        h := hmac.New(sha256.New, []byte(clientSecret))
        h.Write([]byte(username + clientID))
        return base64.StdEncoding.EncodeToString(h.Sum(nil))
    }(clientSecret, clientID, username)
 
    actualHash := ComputeSecretHash(clientSecret, clientID, username)
    assert.Equal(t, expectedHash, actualHash)
}
 