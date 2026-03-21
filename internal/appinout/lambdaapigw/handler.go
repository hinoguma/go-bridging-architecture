package lambdaapigw

import (
	"context"
)

type LambdaAPIGWHandler interface {
	Do(ctx context.Context, request HandlerRequest) (HandlerResponse, error)
}
