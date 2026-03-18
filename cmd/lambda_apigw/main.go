package main

import (
	"app/pkg/crosscutting"
	"app/pkg/crosscutting/errors"
	"app/pkg/crosscutting/infra"
	"app/pkg/crosscutting/log"
	"app/pkg/crosscutting/timer"
	"app/pkg/driven_infra/rdb"
	"app/pkg/driver_interface/lambdaapigw"
	"app/pkg/setup"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func main() {

	// cold start phase: initialize application and infrastructure
	err := coldStart()
	if err != nil {
		fmt.Printf(errors.ToJsonString(err))
		panic(err)
	}

	// start lambda: this is like web server start.
	lambda.Start(lambdaHandler)
}

var apiRouter APIRouter

func coldStart() error {
	// setup application
	ctx := context.Background()
	ctx = context.WithValue(ctx, lambdaapigw.RequestIDKey, "cold_start_request_id")

	// set up crosscutting infrastructure
	crosscutting.SetGlobalRandStrIDGenerator(
		infra.NewUUIDV4Generator(),
	)
	timer.SetGlobalTimeGenerator(infra.NewJstTimeGenerator())
	errors.SetContextRequestIDKey(lambdaapigw.RequestIDKey)

	// set up user call -> app infra
	db, err := rdb.NewPostgresDBFromEnv()
	if err != nil {
		return errors.LiftWithCtx(err, ctx)
	}
	err = setup.GetDrivenInfraRegistry().Initialize(ctx, db)
	if err != nil {
		return errors.LiftWithCtx(err, ctx)
	}
	setup.GetDomainRepositoryRegistry().Initialize(ctx, setup.GetDrivenInfraRegistry())
	setup.GetDomainServiceRegistry().Initialize(
		ctx, setup.GetDomainRepositoryRegistry(),
	)
	setup.GetUseCaseRegistry().Initialize(
		ctx,
		setup.GetDomainRepositoryRegistry(),
		setup.GetDomainServiceRegistry(),
	)
	setup.GetDriverInterfaceRegistry().Initialize(ctx, setup.GetUseCaseRegistry())

	apiRouter = NewAPIRouter(setup.GetDriverInterfaceRegistry())
	return nil
}

func lambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	LambdaContext, ok := lambdacontext.FromContext(ctx)
	if ok {
		ctx = context.WithValue(ctx, lambdaapigw.RequestIDKey, LambdaContext.AwsRequestID)
	} else {
		ctx = context.WithValue(ctx, lambdaapigw.RequestIDKey, "unknown_request_id")
	}

	beforeMiddlewares, handler, afterMiddlewares := apiRouter.Do(ctx, event)

	runner := lambdaapigw.NewLambdaAPIGWHandlerRunner(
		beforeMiddlewares, handler, afterMiddlewares,
	)
	handlerResp, err := runner.Run(ctx, lambdaapigw.NewHandlerRequest(event))
	if err != nil {
		log.Error(log.ErrorLogRequest{Err: err})
	}
	return handlerResp.Raw, err
}
