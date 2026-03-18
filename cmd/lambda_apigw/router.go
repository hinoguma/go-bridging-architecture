package main

import (
	"app/internal/driver_interface/lambdaapigw"
	"app/internal/setup"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type APIRouter struct {
	driverInterfaceRegistry           *setup.DriverInterfaceRegistry
	driverInterfaceMiddlewareRegistry *setup.DriverInterfaceMiddlewareRegistry
}

func (router APIRouter) Do(
	ctx context.Context, event events.APIGatewayProxyRequest,
) (lambdaapigw.BeforeMiddlewareCollection, lambdaapigw.LambdaAPIGWBHandler, lambdaapigw.AfterMiddlewareCollection) {

	beforeMiddlewares := lambdaapigw.BeforeMiddlewareCollection{}
	afterMiddlewares := lambdaapigw.AfterMiddlewareCollection{}
	var handler lambdaapigw.LambdaAPIGWBHandler

	beforeMiddlewares.AddItem(
		router.driverInterfaceMiddlewareRegistry.AuthenticateBeforeMiddleware(),
	)

	if event.Path == "/bank/account/deposit" && event.HTTPMethod == "POST" {
		handler = router.driverInterfaceRegistry.DepositHandler()
	} else if event.Path == "/bank/account/withdraw" && event.HTTPMethod == "POST" {
		handler = router.driverInterfaceRegistry.WithdrawHandler()
	}

	return beforeMiddlewares, handler, afterMiddlewares
}

func NewAPIRouter(
	driverInterfaceRegistry *setup.DriverInterfaceRegistry,
) APIRouter {
	return APIRouter{
		driverInterfaceRegistry: driverInterfaceRegistry,
	}
}
