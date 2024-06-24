package controllers

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"
// 	"time"
	
// "go-chat-app/models"


// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/secretsmanager"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // Mock implementations

// // MockSecretManagerClient implements a mock SecretsManagerClient interface
// type MockSecretManagerClient struct{}

// // GetSecretValue mocks fetching a secret value from Secrets Manager
// func (m *MockSecretManagerClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
// 	mockOutput := &secretsmanager.GetSecretValueOutput{
// 		SecretString: aws.String(`{"UserPoolID":"mock-user-pool-id","Region":"mock-region"}`),
// 	}
// 	if *input.SecretId != "mock-secret" {
// 		return nil, errors.New("secret not found")
// 	}
// 	return mockOutput, nil
// }

// // MockSingleResult simulates a mock mongo.SingleResult
// type MockSingleResult struct {
// 	user models.User
// 	err  error
// }

// // Decode mocks decoding a result into a user
// func (m *MockSingleResult) Decode(v interface{}) error {
// 	userPtr := v.(*models.User)
// 	*userPtr = m.user
// 	return m.err
// }

// // MockFindOneUser mocks finding a user in MongoDB
// func MockFindOneUser(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *MockSingleResult {
// 	// Mock implementation to return a predefined user or error
// 	mockUser := models.User{
// 		First_name: "John",
// 		Last_name:  "Doe",
// 		Password:   "mockpassword",
// 		Email:      "john.doe@example.com",
// 		Phone:      "1234567890",
// 		User_id:    "1234567890",
// 	}
// 	return &MockSingleResult{user: mockUser, err: nil} // Adjust as needed for different scenarios
// }

// func TestGetUser(t *testing.T) {
// 	// Override the actual functions with mock implementations
// 	getSecret = func(client SecretsManagerClient, secretName string) (*models.Secret, error) {
// 		mockClient := &MockSecretManagerClient{}
// 		input := &secretsmanager.GetSecretValueInput{
// 			SecretId: aws.String(secretName),
// 		}
// 		output, err := mockClient.GetSecretValue(context.Background(), input)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var secret models.Secret
// 		err = json.Unmarshal([]byte(*output.SecretString), &secret)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &secret, nil
// 	}
// 	findOneUser = func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
// 		return MockFindOneUser(ctx, filter, opts...)
// 	}

// 	// Create a request to pass to the handler
// 	req, err := http.NewRequest("GET", "/users/1234567890", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a ResponseRecorder to record the response
// 	rr := httptest.NewRecorder()

// 	// Call the handler function directly and pass the Request and ResponseRecorder
// 	GetUser(rr, req)

// 	// Check the status code
// 	if rr.Code != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			rr.Code, http.StatusOK)
// 	}

// 	// Check the response body
// 	expected := models.UserResponse{
// 		FirstName: "John",
// 		LastName:  "Doe",
// 		Password:  "mockpassword", // Note: Normally, you wouldn't expose passwords in responses!
// 		Email:     "john.doe@example.com",
// 		Phone:     "1234567890",
// 		UserID:    "1234567890",
// 	}

// 	var response models.UserResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Errorf("error decoding JSON response: %v", err)
// 	}

// 	if !reflect.DeepEqual(response, expected) {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			response, expected)
// 	}
// }
