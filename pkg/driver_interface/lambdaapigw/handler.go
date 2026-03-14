package lambdaapigw

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaAPIGWBHandler interface {
	Do(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
