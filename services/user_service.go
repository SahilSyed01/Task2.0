package services

import (
	"log"
	"user-management-service/config"
	"user-management-service/models"
	"user-management-service/repository"
	"user-management-service/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func RegisterUser(username, email, password string) (*models.User, error) {
	// Create user in Cognito
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(config.AppConfig.CognitoClientID),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
	_, err := utils.GetCognitoClient().SignUp(input)
	if err != nil {
		log.Println("Error signing up user:", err)
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
	}
	err = repository.CreateUser(user)
	return user, err
}

func LoginUser(username, password string) (string, string, error) {
	// Authenticate user with Cognito
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String("ADMIN_NO_SRP_AUTH"),
		ClientId: aws.String(config.AppConfig.CognitoClientID),
		UserPoolId: aws.String(config.AppConfig.CognitoPoolID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}

	authOutput, err := utils.GetCognitoClient().AdminInitiateAuth(authInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println(aerr.Code(), aerr.Message())
		} else {
			log.Println(err)
		}
		return "", "", err
	}
	cognitoToken := *authOutput.AuthenticationResult.IdToken

	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}

	customToken, err := utils.GenerateJWT(user.ID)
	return cognitoToken, customToken, err
}
