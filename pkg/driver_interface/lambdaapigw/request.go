package lambdaapigw

import "github.com/aws/aws-lambda-go/events"

type HandlerRequest struct {
	Raw events.APIGatewayProxyRequest
}

func NewHandlerRequest(raw events.APIGatewayProxyRequest) HandlerRequest {
	return HandlerRequest{Raw: raw}
}

type HandlerResponse struct {
	Raw events.APIGatewayProxyResponse
}
