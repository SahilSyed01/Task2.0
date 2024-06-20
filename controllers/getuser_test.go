package controllers

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// 	"go-chat-app/models"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // Mocking dependencies
// type MockUserCollection struct {
// 	mock.Mock
// }

// func (m *MockUserCollection) FindOne(ctx context.Context, filter interface{}) SingleResultInterface {
// 	args := m.Called(ctx, filter)
// 	return args.Get(0).(SingleResultInterface)
// }

// type MockSingleResult struct {
// 	mock.Mock
// }

// func (m *MockSingleResult) Decode(v interface{}) error {
// 	args := m.Called(v)
// 	return args.Error(0)
// }


// func TestGetUser(t *testing.T) {
// 	mockUserCollection := new(MockUserCollection)
//     userCollection = mockUserCollection

// 	// Setting environment variables
// 	os.Setenv("REGION", "mockRegion")
// 	os.Setenv("SECRET", "mockSecret")

// 	// Test cases
// 	tests := []struct {
// 		name               string
// 		userID             string
// 		mockFindOneResult  SingleResultInterface
// 		mockFindOneError   error
// 		mockDecodeError    error
// 		expectedStatusCode int
// 		expectedResponse   models.UserResponse
// 	}{
// 		{
// 			name:   "User found",
// 			userID: "123",
// 			mockFindOneResult: func() SingleResultInterface {
// 				result := new(MockSingleResult)
// 				result.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
// 					user := args.Get(0).(*models.User)
// 					user.First_name = "John"
// 					user.Last_name = "Doe"
// 					user.Password = "password"
// 					user.Email = "john.doe@example.com"
// 					user.Phone = "123-456-7890"
// 					user.User_id = "123"
// 				})
// 				return result
// 			}(),
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse: models.UserResponse{
// 				FirstName: "John",
// 				LastName:  "Doe",
// 				Password:  "password",
// 				Email:     "john.doe@example.com",
// 				Phone:     "123-456-7890",
// 				UserID:    "123",
// 			},
// 		},
// 		{
// 			name:               "User not found",
// 			userID:             "456",
// 			mockFindOneResult:  new(MockSingleResult),
// 			mockFindOneError:   mongo.ErrNoDocuments,
// 			expectedStatusCode: http.StatusNotFound,
// 		},
// 		{
// 			name:               "Internal server error",
// 			userID:             "789",
// 			mockFindOneResult:  new(MockSingleResult),
// 			mockFindOneError:   errors.New("some error"),
// 			expectedStatusCode: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockFindOneResult := new(MockSingleResult)
// 			mockFindOneResult.On("Decode", mock.Anything).Return(tt.mockDecodeError)
// 			mockUserCollection.On("FindOne", mock.Anything, bson.M{"user_id": tt.userID}).Return(mockFindOneResult)

// 			// Create a request to pass to the handler
// 			req, err := http.NewRequest("GET", "/users/"+tt.userID, nil)
// 			assert.NoError(t, err)

// 			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 			rr := httptest.NewRecorder()
// 			handler := http.HandlerFunc(GetUser)

// 			// Perform the request
// 			handler.ServeHTTP(rr, req)

// 			// Check the status code is what we expect.
// 			assert.Equal(t, tt.expectedStatusCode, rr.Code)

// 			// Check the response body is what we expect.
// 			if tt.expectedStatusCode == http.StatusOK {
// 				var response models.UserResponse
// 				err := json.NewDecoder(rr.Body).Decode(&response)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedResponse, response)
// 			}
// 		})
// 	}
// }
