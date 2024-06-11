package controllers

import (
    "bytes"
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// Mocking the MongoDB collection
type mockCollection struct {
    data []bson.M
}

func (m *mockCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
    return nil, nil // Dummy implementation
}

func (m *mockCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
    return nil // Dummy implementation
}

func TestSignup(t *testing.T) {
    // Mock the HTTP request and response
    requestBody := []byte(`{"email": "test@example.com", "password": "password"}`)
    req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Authorization", "Bearer eyJraWQiOiJlOTJxZGdvRDFKSVUrZEhoRE9jMnBtbXUzN0JoSmNtTWExQTA4YmFNY1hJPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIyNGE4NDQ5OC1mMDgxLTcwZWQtMTAzMS04NDlkOGVmNDAxZmMiLCJhdWQiOiI1aGkzcDBkMGx2cDdmY2wxbzA1ZmNoajh1aSIsImV2ZW50X2lkIjoiNmI2MmQzMTktNmZmNC00YWRiLTlkM2EtNDRhYTkxMTVjNzg4IiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE3MTgwOTEwNzgsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0xX2ljWGVnMmVpdiIsImNvZ25pdG86dXNlcm5hbWUiOiJteXRlc3R1c2VyIiwiZXhwIjoxNzE4MTc3NDc4LCJpYXQiOjE3MTgwOTEwNzh9.s739RhmqlG-io4rtk9RQrDXF7FirSkkbo8WB7JNwCM8uPMSvJzDqgIhUzzhqFe1OS-Iq_gG-3CBnfO-hUTKBFrmZDCmuhH0CFbiqs4iDNUpgJP-I8Mv7AhpMxp_6nEY0hS1cyiJuWeLXIztxm7l0ogolUROt9kYqC6v26AVu65aoE4RTgku_Yzg6PKIkJbJFiRq0vzfUXmYUuPksFa5nQ_5IBzFa8NrVL6O6qS6r7kQX3rnMTSzFNa5vs5tOMCm61XCYqcAFD3IzWhxXh5O9CMyRZiwdi5kCDNe4NyONwiX3oIKFtL2bfFJuGP0HwiwEkPKAmXY53EoMx8mGgJIZoA")

    rr := httptest.NewRecorder()

    Signup(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

}

func TestLogin(t *testing.T) {
    requestBody := []byte(`{"email": "test@example.com", "password": "password"}`)
    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Authorization", "Bearer your-valid-jwt-token")

    rr := httptest.NewRecorder()

    Login(rr, req)

    // Check the HTTP status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Optionally, you can check the response body or other aspects of the response
}

func TestGetUsers(t *testing.T) {
    // Mock the HTTP request and response
    req, err := http.NewRequest("GET", "/users", nil)
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Authorization", "Bearer your-valid-jwt-token")

    rr := httptest.NewRecorder()

    GetUsers(rr, req)

    // Check the HTTP status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Optionally, you can check the response body or other aspects of the response
}

func TestGetUser(t *testing.T) {
    // Mock the HTTP request and response
    req, err := http.NewRequest("GET", "/users/123", nil)
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Authorization", "Bearer your-valid-jwt-token")

    rr := httptest.NewRecorder()

    GetUser(rr, req)

    // Check the HTTP status code
    assert.Equal(t, http.StatusOK, rr.Code)
}
