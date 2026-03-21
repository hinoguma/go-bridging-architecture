package main

import (
	"app/internal/callapp/lambdaapigw"
	"app/internal/setup"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type APIRouter struct {
	callAppRegistry           *setup.CallAppRegistry
	callAppMiddlewareRegistry *setup.CallAppMiddlewareRegistry
}

func (router APIRouter) Do(
	ctx context.Context, event events.APIGatewayProxyRequest,
) (lambdaapigw.BeforeMiddlewareCollection, lambdaapigw.LambdaAPIGWHandler, lambdaapigw.AfterMiddlewareCollection) {

	beforeMiddlewares := lambdaapigw.BeforeMiddlewareCollection{}
	afterMiddlewares := lambdaapigw.AfterMiddlewareCollection{}
	var handler lambdaapigw.LambdaAPIGWHandler

	beforeMiddlewares.AddItem(
		router.callAppMiddlewareRegistry.AuthenticateBeforeMiddleware(),
	)

	if event.Path == "/bank/account/deposit" && event.HTTPMethod == "POST" {
		handler = router.callAppRegistry.DepositHandler()
	} else if event.Path == "/bank/account/withdraw" && event.HTTPMethod == "POST" {
		handler = router.callAppRegistry.WithdrawHandler()
	} else {

	}

	return beforeMiddlewares, handler, afterMiddlewares
}

func NewAPIRouter(
	callAppRegistry *setup.CallAppRegistry,
	callAppMiddlewareRegistry *setup.CallAppMiddlewareRegistry,
) APIRouter {
	return APIRouter{
		callAppRegistry:           callAppRegistry,
		callAppMiddlewareRegistry: callAppMiddlewareRegistry,
	}
}
