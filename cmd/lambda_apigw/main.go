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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func main() {

	// cold start phase: initialize application and infrastructure
	err := coldStart()
	if err != nil {
		log.Error(
			"error in cold start",
			log.WithErr(err),
		)
		panic(err)
	}

	// start lambda: this is like web server start.
	lambda.Start(lambdaHandler)
}

var apiRouter APIRouter

func coldStart() error {
	log.Info("cold start started")

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

	log.Info("cold start completed")
	return nil
}

func lambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestID := "unknown_request_id"
	LambdaContext, ok := lambdacontext.FromContext(ctx)
	if ok {
		requestID = LambdaContext.AwsRequestID
	}
	ctx = context.WithValue(ctx, lambdaapigw.RequestIDKey, requestID)

	log.Info(
		"lambda handler started",
		log.WithTags(map[string]any{"event": event}),
		log.WithRequestID(requestID),
	)

	beforeMiddlewares, handler, afterMiddlewares := apiRouter.Do(ctx, event)

	runner := lambdaapigw.NewLambdaAPIGWHandlerRunner(
		beforeMiddlewares, handler, afterMiddlewares,
	)
	handlerResp, err := runner.Run(ctx, lambdaapigw.NewHandlerRequest(event))
	if err != nil {
		log.Error(
			"error in lambda handler",
			log.WithErr(err),
			log.WithRequestID(requestID),
		)
	}
	return handlerResp.Raw, err
}
