// helpers/tokenHelper_test.go

package helpers

import (
	"testing"

	// "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking jwt.Token for testing
type MockToken struct {
    mock.Mock
}

func (m *MockToken) SignedString(key []byte) (string, error) {
    args := m.Called(key)
    return args.String(0), args.Error(1)
}

func TestGenerateToken(t *testing.T) {
    // Mocking SECRET_KEY
    SECRET_KEY = "secret"

    // Test positive case
    firstName := "John"
    userID := "123"
    expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInVpZCI6IjEyMyJ9.wKBasWto25zYrg_X_WuMcJde1TJ4RB5EyYqJmw-dBC4"

    // Mock token
    mockToken := new(MockToken)
    mockToken.On("SignedString", []byte("secret")).Return(expectedToken, nil)

    // Call the method being tested
    token, err := GenerateToken(firstName, userID)

    // Assert the result
    if err != nil {
        t.Errorf("Error generating token: %v", err)
    }
    if token != expectedToken {
        t.Errorf("Generated token does not match expected token. Expected: %s, Got: %s", expectedToken, token)
    }
}


func TestValidateToken(t *testing.T) {
    // Mocking SECRET_KEY
    SECRET_KEY = "secret"

    // Sample token
    signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInVpZCI6IjEyMyJ9.wKBasWto25zYrg_X_WuMcJde1TJ4RB5EyYqJmw-dBC4"

    // Call the method being tested
    claims, err := ValidateToken(signedToken)

    // Assert the result
    assert.Nil(t, err)
    assert.Equal(t, "John", claims.First_name)
    assert.Equal(t, "123", claims.Uid)

    // Test negative case (error in parsing token)
    signedToken = "invalid-token"
    _, err = ValidateToken(signedToken)

    // Assert the result
    assert.NotNil(t, err)
}


func TestValidateToken_ErrorCastingClaims(t *testing.T) {
    // Mock token claims
    invalidClaims := &SignedDetails{}

    // Mock token to return invalid claims
    mockToken := new(MockToken)
    mockToken.On("Claims").Return(invalidClaims)

    // Call the method being tested
    _, err := ValidateToken("valid-signed-token")

    // Assertions
    assert.Error(t, err)
}


func TestValidateToken_EmptyToken(t *testing.T) {
    // Call the method being tested with an empty token
    _, err := ValidateToken("")

    // Assertions
    assert.Error(t, err)
}
