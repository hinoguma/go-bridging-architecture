package setup

import (
	"app/internal/driver_entrance/lambdaapigw"
	"app/internal/driver_entrance_applogic_bridge/lambdaapigw_bridge"
	"context"
	"sync"
)

type DriverInterfaceMiddlewareRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry

	authenticateBeforeMiddleware lambdaapigw.BeforeMiddleware
}

func (registry *DriverInterfaceMiddlewareRegistry) Initialize(ctx context.Context, useCaseRegistry *UseCaseRegistry) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry
	registry.authenticateBeforeMiddleware = lambdaapigw_bridge.NewAuthenticateBeforeMiddleware(
		registry.useCaseRegistry.BankAccountAuthMWUseCase(),
	)
	registry.Unlock()
}

func (registry *DriverInterfaceMiddlewareRegistry) AuthenticateBeforeMiddleware() lambdaapigw.BeforeMiddleware {
	if registry.authenticateBeforeMiddleware == nil {
		registry.Lock()
		registry.authenticateBeforeMiddleware = lambdaapigw_bridge.NewAuthenticateBeforeMiddleware(
			registry.useCaseRegistry.BankAccountAuthMWUseCase(),
		)
		registry.Unlock()
	}
	return registry.authenticateBeforeMiddleware
}

var driverInterfaceMiddlewareRegistry DriverInterfaceMiddlewareRegistry = DriverInterfaceMiddlewareRegistry{}

func GetDriverInterfaceMiddlewareRegistry() *DriverInterfaceMiddlewareRegistry {
	return &driverInterfaceMiddlewareRegistry
}
