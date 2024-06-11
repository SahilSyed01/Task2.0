package helpers_test

import (
	"errors"
	"go-chat-app/helpers"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

type mockToken struct {
	tokenString string
	err         error
}

func (m *mockToken) SignedString(_ []byte) (string, error) {
	return m.tokenString, m.err
}

func (m *mockToken) ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return &jwt.Token{}, m.err
}

func TestGenerateToken(t *testing.T) {
	testCases := []struct {
		name       string
		firstName  string
		userID     string
		secretKey  string
		mockToken  *mockToken
		expectErr  bool
	}{
		{
			name:      "Successful token generation",
			firstName: "John",
			userID:    "123456",
			secretKey: "secret",
			mockToken: &mockToken{
				tokenString: "generated_token",
				err:         nil,
			},
			expectErr: false,
		},
		// Add more test cases here as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helpers.SECRET_KEY = tc.secretKey
			helpers.MockTokenGenerator = tc.mockToken
			_, err := helpers.GenerateToken(tc.firstName, tc.userID)
			assert.Equal(t, tc.expectErr, err != nil)
		})
	}
}

func TestValidateToken(t *testing.T) {
    testCases := []struct {
        name      string
        token     string
        secretKey string
        mockToken *mockToken
        expectErr bool
    }{
        {
            name:      "valid token",
            token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
            mockToken: &mockToken{
                tokenString: "valid_token",
                err:         nil,
            },
            expectErr: false,
        },
        {
            name:      "Invalid token",
            token:     "invalid_token",
            mockToken: &mockToken{
                tokenString: "invalid_token",
                err:         errors.New("token validation failed"),
            },
            expectErr: true,
        },
        // Add more test cases here as needed
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            helpers.SECRET_KEY = tc.secretKey
            helpers.MockTokenParser = tc.mockToken
            _, err := helpers.ValidateToken(tc.token)
            assert.Equal(t, tc.expectErr, err != nil)
        })
    }
}
