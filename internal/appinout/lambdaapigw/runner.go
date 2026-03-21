package lambdaapigw

import (
	"app/internal/crosscutting/errors"
	"context"
)

func NewLambdaAPIGWHandlerRunner(
	beforeMiddlewares BeforeMiddlewareCollection,
	handler LambdaAPIGWHandler,
	afterMiddlewares AfterMiddlewareCollection,
) LambdaAPIGWHandlerRunner {
	return LambdaAPIGWHandlerRunner{
		beforeMiddlewares: beforeMiddlewares,
		handler:           handler,
		afterMiddlewares:  afterMiddlewares,
	}
}

type LambdaAPIGWHandlerRunner struct {
	beforeMiddlewares BeforeMiddlewareCollection
	handler           LambdaAPIGWHandler
	afterMiddlewares  AfterMiddlewareCollection
}

func (runner LambdaAPIGWHandlerRunner) Run(ctx context.Context, request HandlerRequest) (HandlerResponse, error) {
	for _, beforeMW := range runner.beforeMiddlewares.GetItems() {
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
		return NewInternalServerErrorResponse(), err
	}

	for _, afterMW := range runner.afterMiddlewares.GetItems() {
		afterMWRes, err := afterMW.Do(ctx, request, resp)
		if err != nil {
			// if afterMW returns error, we return the original response and log the error
			return resp, errors.Lift(err)
		}
		resp = afterMWRes.Response
	}

	return resp, nil
}
