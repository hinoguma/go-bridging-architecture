package main

import (
	"app/internal/appinfra/postgres"
	"app/internal/appinout/lambdaapigw"
	"app/internal/crosscutting"
	"app/internal/crosscutting/errors"
	"app/internal/crosscutting/infra"
	"app/internal/crosscutting/log"
	"app/internal/crosscutting/timer"
	"app/internal/setup"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func main() {

	// lambda cold start phase
	app := application{}
	err := app.coldStart()
	if err != nil {
		log.Error(
			"error in cold start",
			log.WithErr(err),
		)
		panic(err)
	}

	// start lambda: this is like web server start.
	lambda.Start(app.lambdaHandler)
}

var apiRouter lambdaapigw.APIRouter

type application struct {
	router *lambdaapigw.APIRouter
}

func (app *application) coldStart() error {
	log.Info("cold start started")

	ctx := context.Background()
	ctx = context.WithValue(ctx, lambdaapigw.RequestIDKey, "cold_start_request_id")

	// set up crosscutting infrastructure
	crosscutting.SetGlobalRandStrIDGenerator(
		infra.NewUUIDV4Generator(),
	)
	timer.SetGlobalTimeGenerator(infra.NewJstTimeGenerator())
	errors.SetContextRequestIDKey(lambdaapigw.RequestIDKey)

	db, err := postgres.NewPostgresDBFromEnv()
	if err != nil {
		return errors.LiftWithCtx(err, ctx)
	}
	err = setup.GetAppInfraRegistry().Initialize(ctx, db)
	if err != nil {
		return errors.LiftWithCtx(err, ctx)
	}
	setup.GetDomainRepositoryRegistry().Initialize(ctx, setup.GetAppInfraRegistry())
	setup.GetDomainServiceRegistry().Initialize(
		ctx, setup.GetDomainRepositoryRegistry(),
	)
	setup.GetUseCaseRegistry().Initialize(
		ctx,
		setup.GetDomainRepositoryRegistry(),
		setup.GetDomainServiceRegistry(),
	)
	setup.GetAppInOutMiddlewareRegistry().Initialize(ctx, setup.GetUseCaseRegistry())
	setup.GetAppInOutRegistry().Initialize(ctx, setup.GetUseCaseRegistry())

	router := lambdaapigw.NewAPIRouter(
		setup.GetAppInOutRegistry(),
		setup.GetAppInOutMiddlewareRegistry(),
	)
	app.router = &router

	log.Info("cold start completed")
	return nil
}

func (app *application) lambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	if app.router == nil {
		err := errors.NewWithCtx("router is not initialized", ctx)
		log.Error(
			"router is not initialized",
			log.WithErr(err),
			log.WithRequestID(requestID),
		)
		return lambdaapigw.NewInternalServerErrorResponse().Raw, err
	}

	beforeMiddlewares, handler, afterMiddlewares := app.router.Do(ctx, event)
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
