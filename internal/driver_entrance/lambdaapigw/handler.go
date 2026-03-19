package lambdaapigw

import (
	"context"
)

type LambdaAPIGWBHandler interface {
	Do(ctx context.Context, request HandlerRequest) (HandlerResponse, error)
}
