package lambdaapigw

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type APIRouter struct {
	handlerRegistry    HandlerRegistry
	middlewareRegistry MiddlewareRegistry
}

func (router APIRouter) Do(
	ctx context.Context, event events.APIGatewayProxyRequest,
) (BeforeMiddlewareCollection, LambdaAPIGWHandler, AfterMiddlewareCollection) {

	beforeMiddlewares := BeforeMiddlewareCollection{}
	afterMiddlewares := AfterMiddlewareCollection{}
	var handler LambdaAPIGWHandler

	beforeMiddlewares.AddItem(
		router.middlewareRegistry.AuthenticateBeforeMiddleware(),
	)

	if event.Path == "/bank/account/deposit" && event.HTTPMethod == "POST" {
		handler = router.handlerRegistry.DepositHandler()

	} else if event.Path == "/bank/account/withdraw" && event.HTTPMethod == "POST" {
		handler = router.handlerRegistry.WithdrawHandler()
	} else {

	}

	return beforeMiddlewares, handler, afterMiddlewares
}

func NewAPIRouter(
	handlerRegistry HandlerRegistry,
	middlewareRegistry MiddlewareRegistry,
) APIRouter {
	return APIRouter{
		handlerRegistry:    handlerRegistry,
		middlewareRegistry: middlewareRegistry,
	}
}
