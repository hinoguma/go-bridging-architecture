package lambdaapigw

import "context"

type BeforeMiddlewareResult struct {
	Ctx      context.Context
	Request  HandlerRequest
	Response *HandlerResponse
}

type BeforeMiddleware interface {
	Do(ctx context.Context, request HandlerRequest) (BeforeMiddlewareResult, error)
}

type AfterMiddlewareResult struct {
	Ctx      context.Context
	Response HandlerResponse
}

type AfterMiddleware interface {
	Do(
		ctx context.Context,
		request HandlerRequest,
		response HandlerResponse,
	) (AfterMiddlewareResult, error)
}
