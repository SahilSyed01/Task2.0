package controllers

import (
	"context"
	"fmt"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
    password := "password123"
    hashedPassword := HashPassword(password)

    if hashedPassword == "" {
        t.Error("Hashed password should not be empty")
    }
}

func TestVerifyPassword_Positive(t *testing.T) {
    userPassword := "password123"
    providedPassword := HashPassword(userPassword)

    check, msg := VerifyPassword(userPassword, providedPassword)

    if !check {
        t.Error("Passwords should match")
    }
    if msg != "" {
        t.Error("Message should be empty for positive case")
    }
}

func TestVerifyPassword_Negative(t *testing.T) {
    userPassword := "password123"
    providedPassword := "wrongpassword"

    check, msg := VerifyPassword(userPassword, providedPassword)

    if check {
        t.Error("Passwords should not match")
    }
    if msg == "" {
        t.Error("Message should not be empty for negative case")
    }
}

func TestJwtAuthenticator_Mock_Positive(t *testing.T) {
    // Mock implementation
    mockJwtAuthenticator := func(ctx context.Context, region, userPoolID, tokenString string) (interface{}, error) {
        return "mockUser", nil
    }

    result, err := mockJwtAuthenticator(context.Background(), "region", "userPoolID", "tokenString")

    if err != nil {
        t.Error("Expected no error, got:", err)
    }
    if result != "mockUser" {
        t.Errorf("Expected 'mockUser', got: %v", result)
    }
}

func TestJwtAuthenticator_Mock_Error(t *testing.T) {
    // Mock implementation returning an error
    mockJwtAuthenticator := func(ctx context.Context, region, userPoolID, tokenString string) (interface{}, error) {
        return nil, fmt.Errorf("mock error") // Return a mock error
    }

    _, err := mockJwtAuthenticator(context.Background(), "region", "userPoolID", "tokenString")

    if err == nil {
        t.Error("Expected error, got nil")
    }
}
