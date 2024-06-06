package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	//"user-management-service/config"
)

var cognitoClient *cognitoidentityprovider.CognitoIdentityProvider

func InitCognito() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	cognitoClient = cognitoidentityprovider.New(sess)
}

func GetCognitoClient() *cognitoidentityprovider.CognitoIdentityProvider {
	return cognitoClient
}
