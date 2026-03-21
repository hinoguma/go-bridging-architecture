package setup

import (
	"app/internal/callapp/lambdaapigw"
	"app/internal/callapp_applogic_connector/lambdaapigw_connector"
	"context"
	"sync"
)

type CallAppMiddlewareRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry

	authenticateBeforeMiddleware lambdaapigw.BeforeMiddleware
}

func (registry *CallAppMiddlewareRegistry) Initialize(ctx context.Context, useCaseRegistry *UseCaseRegistry) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry
	registry.authenticateBeforeMiddleware = lambdaapigw_connector.NewAuthenticateBeforeMiddleware(
		registry.useCaseRegistry.BankAccountAuthMWUseCase(),
	)
	registry.Unlock()
}

func (registry *CallAppMiddlewareRegistry) AuthenticateBeforeMiddleware() lambdaapigw.BeforeMiddleware {
	if registry.authenticateBeforeMiddleware == nil {
		registry.Lock()
		registry.authenticateBeforeMiddleware = lambdaapigw_connector.NewAuthenticateBeforeMiddleware(
			registry.useCaseRegistry.BankAccountAuthMWUseCase(),
		)
		registry.Unlock()
	}
	return registry.authenticateBeforeMiddleware
}

var callAppMiddlewareRegistry CallAppMiddlewareRegistry = CallAppMiddlewareRegistry{}

func GetCallAppMiddlewareRegistry() *CallAppMiddlewareRegistry {
	return &callAppMiddlewareRegistry
}
