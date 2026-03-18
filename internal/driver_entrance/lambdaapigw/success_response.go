package lambdaapigw

import (
	"app/internal/crosscutting"

	"github.com/aws/aws-lambda-go/events"
)

func NewSuccessResponse(body string) HandlerResponse {
	return HandlerResponse{
		Raw: events.APIGatewayProxyResponse{
			StatusCode: crosscutting.APIStatusOK.Int(),
			Body:       body,
		},
	}
}

func NewAcceptedResponse(body string) HandlerResponse {
	return HandlerResponse{
		Raw: events.APIGatewayProxyResponse{
			StatusCode: crosscutting.APIStatusAccepted.Int(),
			Body:       body,
		},
	}
}
