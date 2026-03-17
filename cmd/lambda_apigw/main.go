package main

import (
	"app/pkg/crosscutting"
	"app/pkg/crosscutting/errors"
	"app/pkg/crosscutting/infra"
	"app/pkg/crosscutting/timer"
	"app/pkg/driver_interface/lambdaapigw"
	"app/pkg/setup"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	// cold start
	err := coldStart()
	if err != nil {
		fmt.Printf(errors.ToJsonString(err))
		panic(err)
	}

	// start lambda
	lambda.Start(lambdaHandler)
}

var apiRouter APIRouter

func coldStart() error {
	// setup application
	ctx := context.Background()

	// set up crosscutting infrastructure
	crosscutting.SetGlobalRandStrIDGenerator(
		infra.NewUUIDV4Generator(),
	)
	timer.SetGlobalTimeGenerator(infra.NewJstTimeGenerator())

	// set up user call -> app infra
	//err := setup.GetAppInfraRegistry().Initialize(ctx)
	//if err != nil {
	//	return errors.LiftWithCtx(err, ctx)
	//}
	//setup.GetAppLogicToAppInfraBridgeRegistry().Initialize(ctx, setup.GetAppInfraRegistry())
	//setup.GetAppLogicServiceRegistry().Initialize(ctx, setup.GetAppLogicToAppInfraBridgeRegistry())
	//setup.GetAppLogicUseCaseRegistry().Initialize(ctx, setup.GetAppLogicToAppInfraBridgeRegistry(), setup.GetAppLogicServiceRegistry())
	setup.GetDriverInterfaceRegistry().Initialize(ctx, setup.GetUseCaseRegistry())

	apiRouter = NewAPIRouter(*setup.GetDriverInterfaceRegistry())
	return nil
}

func lambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	beforeMiddlewares, handler, afterMiddlewares := apiRouter.Do(ctx, event)

	runner := lambdaapigw.NewLambdaAPIGWHandlerRunner(
		beforeMiddlewares, handler, afterMiddlewares,
	)
	handlerResp, err := runner.Run(ctx, lambdaapigw.NewHandlerRequest(event))
	if err != nil {
		fmt.Printf(errors.ToJsonString(err))
	}

	return handlerResp.Raw, err
}
