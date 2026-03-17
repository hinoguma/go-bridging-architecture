package main

import (
	"app/pkg/driver_interface/lambdaapigw"
	"app/pkg/setup"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type APIRouter struct {
	driverInterfaceRegistry setup.DriverInterfaceRegistry
}

func (router APIRouter) Do(
	ctx context.Context, event events.APIGatewayProxyRequest,
) ([]lambdaapigw.BeforeMiddleware, lambdaapigw.LambdaAPIGWBHandler, []lambdaapigw.AfterMiddleware) {

	beforeMiddlewares := make([]lambdaapigw.BeforeMiddleware, 0)
	afterMiddlewares := make([]lambdaapigw.AfterMiddleware, 0)
	var handler lambdaapigw.LambdaAPIGWBHandler

	if event.Path == "/bank/account/deposit" && event.HTTPMethod == "POST" {
		handler = router.driverInterfaceRegistry.GetDepositHandler()
	} else if event.Path == "/bank/account/withdraw" && event.HTTPMethod == "POST" {
		// handler = router.driverInterfaceRegistry.GetWithdrawHandler()
	}

	return beforeMiddlewares, handler, afterMiddlewares
}

func NewAPIRouter(
	driverInterfaceRegistry setup.DriverInterfaceRegistry,
) APIRouter {
	return APIRouter{
		driverInterfaceRegistry: driverInterfaceRegistry,
	}
}
