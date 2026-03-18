package main

import (
	"app/internal/crosscutting"
	"app/internal/crosscutting/errors"
	"app/internal/crosscutting/infra"
	"app/internal/crosscutting/log"
	"app/internal/crosscutting/timer"
	"app/internal/driven_infra/postgres"
	"app/internal/driver_entrance/lambdaapigw"
	"app/internal/setup"
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
	db, err := postgres.NewPostgresDBFromEnv()
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
	setup.GetDriverInterfaceMiddlewareRegistry().Initialize(ctx, setup.GetUseCaseRegistry())
	setup.GetDriverInterfaceRegistry().Initialize(ctx, setup.GetUseCaseRegistry())

	apiRouter = NewAPIRouter(
		setup.GetDriverInterfaceRegistry(),
		setup.GetDriverInterfaceMiddlewareRegistry(),
	)

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
	if handler == nil {
		handlerResp := lambdaapigw.NewNotFoundErrorResponse()
		return handlerResp.Raw, nil
	}

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
