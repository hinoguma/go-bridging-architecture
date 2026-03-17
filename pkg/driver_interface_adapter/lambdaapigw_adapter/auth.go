package lambdaapigw_adapter

import (
	"app/pkg/driver_interface/lambdaapigw"
	"context"
)

func NewAuthenticateBeforeMiddleware() lambdaapigw.BeforeMiddleware {
	return authenticateBeforeMiddleware{}
}

type authenticateBeforeMiddleware struct {
}

func (mw authenticateBeforeMiddleware) Do(
	ctx context.Context,
	request lambdaapigw.HandlerRequest,
) (lambdaapigw.BeforeMiddlewareResult, error) {

}
