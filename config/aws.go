package config

import "os"

func GetCognitoClientID() string {
	return os.Getenv("COGNITO_CLIENT_ID")
}

func GetCognitoUserPoolID() string {
	return os.Getenv("COGNITO_USER_POOL_ID")
}
