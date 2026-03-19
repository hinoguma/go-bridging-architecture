package lambdaapigw

import "github.com/aws/aws-lambda-go/events"

type HandlerResponse struct {
	Raw events.APIGatewayProxyResponse
}
