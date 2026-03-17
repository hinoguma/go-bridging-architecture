package lambdaapigw

import (
	"app/pkg/crosscutting/errors"
	"context"
)

type LambdaAPIGWBHandler interface {
	Do(ctx context.Context, request HandlerRequest) (HandlerResponse, error)
}

func NewLambdaAPIGWHandlerRunner(
	beforeMiddlewares []BeforeMiddleware,
	handler LambdaAPIGWBHandler,
	afterMiddlewares []AfterMiddleware,
) LambdaAPIGWHandlerRunner {
	return LambdaAPIGWHandlerRunner{
		beforeMiddlewares: beforeMiddlewares,
		handler:           handler,
		afterMiddlewares:  afterMiddlewares,
	}
}

type LambdaAPIGWHandlerRunner struct {
	beforeMiddlewares []BeforeMiddleware
	handler           LambdaAPIGWBHandler
	afterMiddlewares  []AfterMiddleware
}

func (runner LambdaAPIGWHandlerRunner) Run(ctx context.Context, request HandlerRequest) (HandlerResponse, error) {
	for _, beforeMW := range runner.beforeMiddlewares {
		beforeMWRes, err := beforeMW.Do(ctx, request)
		if err != nil {
			return NewInternalServerErrorResponse(), errors.Lift(err)
		}
		if beforeMWRes.Response != nil {
			return *beforeMWRes.Response, nil
		}
		request = beforeMWRes.Request
	}

	resp, err := runner.handler.Do(ctx, request)
	if err != nil {
		return HandlerResponse{}, err
	}

	for _, afterMW := range runner.afterMiddlewares {
		afterMWRes, err := afterMW.Do(ctx, request, resp)
		if err != nil {
			// if afterMW returns error, we return the original response and log the error
			return resp, errors.Lift(err)
		}
		resp = afterMWRes.Response
	}

	return resp, nil
}
