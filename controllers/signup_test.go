package controllers

// import (
//     "bytes"
//     "context"
//     "encoding/json"
//     "net/http"
//     "net/http/httptest"
//     "testing"
//     "time"

//     "go-chat-app/models"

//     "github.com/stretchr/testify/assert"
// )

// func TestSignup_Success(t *testing.T) {
//     // Prepare a request body
//     requestBody := []byte(`{"email": "test@example.com", "password": "password123", "phone": "1234567890"}`)
//     req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
//     if err != nil {
//         t.Fatal(err)
//     }

//     // Create a ResponseRecorder to record the response
//     rr := httptest.NewRecorder()

//     // Mock dependencies
//     authenticate = MockAuthenticate
//     getAWSConfig = MockGetAWSConfig
//     getSMClient = MockGetSMClient
//     GetSecret = MockGetSecret
//     countDocs = MockCountDocs
//     hashPassword = MockHashPassword
//     insertOneUser = MockInsertOneUser

//     // Call the Signup handler function directly
//     Signup(rr, req)

//     // Check the status code
//     if status := rr.Code; status != http.StatusOK {
//         t.Errorf("handler returned wrong status code: got %v want %v",
//             status, http.StatusOK)
//     }

//     // Check the response body if needed (depends on your implementation)
//     // For example, you might decode and verify the response JSON:
//     var response map[string]interface{}
//     if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
//         t.Errorf("error decoding JSON response: %v", err)
//     }

//     // Assert expected response if needed
//     // Example: assert.Equal(t, expectedResponse, response)
// }

// func TestSignup_BadRequestInvalidData(t *testing.T) {
//     // Prepare a request body with invalid data (missing required fields, etc.)
//     requestBody := []byte(`{"email": "test@example.com"}`)
//     req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
//     if err != nil {
//         t.Fatal(err)
//     }

//     // Create a ResponseRecorder to record the response
//     rr := httptest.NewRecorder()

//     // Mock dependencies
//     authenticate = MockAuthenticate
//     getAWSConfig = MockGetAWSConfig
//     getSMClient = MockGetSMClient
//     GetSecret = MockGetSecret
//     countDocs = MockCountDocs
//     hashPassword = MockHashPassword
//     insertOneUser = MockInsertOneUser

//     // Call the Signup handler function directly
//     Signup(rr, req)

//     // Check the status code
//     if status := rr.Code; status != http.StatusBadRequest {
//         t.Errorf("handler returned wrong status code: got %v want %v",
//             status, http.StatusBadRequest)
//     }

//     // Check the response body if needed
//     // Example: assert.Contains(t, rr.Body.String(), "validation error message")
// }

// // Mocks for dependencies

// func MockAuthenticate(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Simulate authentication middleware behavior
//         // For simplicity, let's assume authentication always succeeds.
//         ctx := context.WithValue(r.Context(), "authenticated", true)
//         next.ServeHTTP(w, r.WithContext(ctx))
//     })
// }

// func MockGetAWSConfig() (string, error) {
//     // Simple mock AWS configuration retrieval
//     return "mocked_aws_config", nil
// }

// func MockGetSMClient(awsConfig string) interface{} {
//     // Simple mock Secrets Manager client
//     return nil // Replace with appropriate mock behavior if needed
// }

// func MockGetSecret(smClient interface{}, secretName string) (*models.Secret, error) {
//     // Simple mock secret retrieval
//     return &models.Secret{
//         UserPoolID: "mock_user_pool_id",
//         Region:     "mock_region",
//     }, nil
// }

// func MockCountDocs(ctx context.Context, filter interface{}) (int64, error) {
//     // Mock count documents function
//     // Return a predefined count for testing scenarios
//     return 0, nil // Replace with appropriate mock behavior
// }

// func MockHashPassword(password string) string {
//     // Simple mock hash password function
//     return "mocked_hashed_password"
// }

// func MockInsertOneUser(ctx context.Context, document interface{}) (*models.InsertOneResult, error) {
//     // Mock insert one user function
//     // Return a predefined result for testing scenarios
//     return &models.InsertOneResult{
//         InsertedID: "mocked_inserted_id",
//     }, nil
// }
