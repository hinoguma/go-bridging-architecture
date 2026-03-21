package setup

import (
	"app/internal/appinout/lambdaapigw"
	"app/internal/appinout_applogic_connector/lambdaapigw_connector"
	"context"
	"sync"
)

type AppInOutMiddlewareRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry

	authenticateBeforeMiddleware lambdaapigw.BeforeMiddleware
}

func (registry *AppInOutMiddlewareRegistry) Initialize(ctx context.Context, useCaseRegistry *UseCaseRegistry) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry
	registry.authenticateBeforeMiddleware = lambdaapigw_connector.NewAuthenticateBeforeMiddleware(
		registry.useCaseRegistry.BankAccountAuthMWUseCase(),
	)
	registry.Unlock()
}

func (registry *AppInOutMiddlewareRegistry) AuthenticateBeforeMiddleware() lambdaapigw.BeforeMiddleware {
	if registry.authenticateBeforeMiddleware == nil {
		registry.Lock()
		registry.authenticateBeforeMiddleware = lambdaapigw_connector.NewAuthenticateBeforeMiddleware(
			registry.useCaseRegistry.BankAccountAuthMWUseCase(),
		)
		registry.Unlock()
	}
	return registry.authenticateBeforeMiddleware
}

var appInOutMiddlewareRegistry AppInOutMiddlewareRegistry = AppInOutMiddlewareRegistry{}

func GetAppInOutMiddlewareRegistry() *AppInOutMiddlewareRegistry {
	return &appInOutMiddlewareRegistry
}
