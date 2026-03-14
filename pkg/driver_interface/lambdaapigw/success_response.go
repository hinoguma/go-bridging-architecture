package lambdaapigw

import (
	"app/pkg/crosscutting"

	"github.com/aws/aws-lambda-go/events"
)

func NewSuccessResponse(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: crosscutting.APIStatusOK.Int(),
		Body:       body,
	}
}

func NewAcceptedResponse(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: crosscutting.APIStatusAccepted.Int(),
		Body:       body,
	}
}
