package helpers
 
import (
    "errors"
    "testing"
 
    "github.com/dgrijalva/jwt-go"
    "github.com/stretchr/testify/assert"
)
 
func TestGenerateToken(t *testing.T) {
    // Mocking SECRET_KEY
    SetSecretKey("secret")
 
    // Test positive case
    firstName := "John"
    userID := "123"
    expectedToken := "expected-token"
 
    // Mock the signing method
    SetSignTokenFunc(func(token *jwt.Token, key interface{}) (string, error) {
        assert.Equal(t, []byte("secret"), key)
 
        // Assert and return expected token
        return expectedToken, nil
    })
 
    // Call the method being tested
    token, err := GenerateToken(firstName, userID)
 
    // Assert the result
    assert.NoError(t, err)
    assert.Equal(t, expectedToken, token)
}
 
func TestGenerateToken_Error(t *testing.T) {
    // Mocking SECRET_KEY
    SetSecretKey("secret")
 
    // Test error case
    firstName := "John"
    userID := "123"
    expectedError := errors.New("error generating token")
 
    // Mock the signing method to return an error
    SetSignTokenFunc(func(token *jwt.Token, key interface{}) (string, error) {
        assert.Equal(t, []byte("secret"), key)
        return "", expectedError
    })
 
    // Call the method being tested
    token, err := GenerateToken(firstName, userID)
 
    // Assert the result
    assert.Empty(t, token)
    assert.Equal(t, expectedError, err)
}
 
func TestValidateToken(t *testing.T) {
    t.Run("test for validate token", func(t *testing.T) {
        // Mocking SECRET_KEY
        SetSecretKey("secret")
 
        // Sample token
        signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInVpZCI6IjEyMyJ9.wKBasWto25zYrg_X_WuMcJde1TJ4RB5EyYqJmw-dBC4"
 
        jwtParseWithClaim = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
            // Assert the keyFunc is correctly passed the secret key
            key, _ := keyFunc(nil)
            assert.Equal(t, []byte("secret"), key)
 
            return &jwt.Token{
                Method: jwt.SigningMethodHS256,
                Raw:    "test",
                Claims: &SignedDetails{
                    First_name:     "John",
                    Uid:            "123",
                    StandardClaims: jwt.StandardClaims{},
                },
                Valid:     true,
                Signature: "test",
            }, nil
        }
 
        // Call the method being tested
        claims, err := ValidateToken(signedToken)
 
        // Assert the result
        assert.NoError(t, err)
        assert.NotNil(t, claims)
        assert.Equal(t, "John", claims.First_name)
        assert.Equal(t, "123", claims.Uid)
    })
 
    t.Run("test for error", func(t *testing.T) {
        SetSecretKey("secret")
        jwtParseWithClaim = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
            return nil, assert.AnError
        }
 
        signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInVpZCI6IjEyMyJ9.wKBasWto25zYrg_X_WuMcJde1TJ4RB5EyYqJmw-dBC4"
 
        claims, err := ValidateToken(signedToken)
        assert.NotNil(t, err)
        assert.Nil(t, claims)
    })
}
 
func TestValidateToken_EmptyToken(t *testing.T) {
    SetSecretKey("secret")
    _, err := ValidateToken("")
 
    assert.Error(t, err)
}
 
func TestSetNewWithClaimsFunc(t *testing.T) {
    // Define a mock function for newWithClaimsFunc
    mockFunc := func(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
        return jwt.New(method)
    }
 
    // Call the function being tested
    SetNewWithClaimsFunc(mockFunc)
 
    // Now, newWithClaimsFunc should be updated with mockFunc
    assert.Equal(t, newWithClaimsFunc(jwt.SigningMethodHS256, jwt.MapClaims{}).Method.Alg(), "HS256")
}
func TestValidateToken_ErrorCastingClaims(t *testing.T) {
    // Mocking SECRET_KEY
    SetSecretKey("secret")
 
    // Sample token
    signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInVpZCI6IjEyMyJ9.wKBasWto25zYrg_X_WuMcJde1TJ4RB5EyYqJmw-dBC4"
 
    jwtParseWithClaim = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
        // Assert the keyFunc is correctly passed the secret key
        key, _ := keyFunc(nil)
        assert.Equal(t, []byte("secret"), key)
 
        // Return a token with incorrect claims type
        return &jwt.Token{
            Method: jwt.SigningMethodHS256,
            Raw:    "test",
            Claims: &jwt.StandardClaims{},
            Valid:  true,
        }, nil
    }
 
    // Call the method being tested
    claims, _ := ValidateToken(signedToken)
 
    assert.Nil(t, claims)
}