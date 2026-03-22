package driven_infra_aws

import (
	"app/config"
)

func GetCognitoClientID() string {
	return config.GetCognitoClientID()
}

func GetCognitoUserPoolID() string {
	return config.GetCognitoUserPoolID()
}
