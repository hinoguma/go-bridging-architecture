package lambdaapigw

import (
	"app/pkg/crosscutting"
	"context"
)

type BeforeMiddlewareResult struct {
	Ctx      context.Context
	Request  HandlerRequest
	Response *HandlerResponse
}

func (result *BeforeMiddlewareResult) SetResponse(response HandlerResponse) {
	result.Response = &response
}

type BeforeMiddleware interface {
	Do(ctx context.Context, request HandlerRequest) (BeforeMiddlewareResult, error)
}

type BeforeMiddlewareCollection struct {
	crosscutting.Collection[BeforeMiddleware]
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

type AfterMiddlewareCollection struct {
	crosscutting.Collection[AfterMiddleware]
}
