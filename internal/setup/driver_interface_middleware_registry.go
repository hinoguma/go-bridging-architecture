package setup

import (
	"app/internal/driver_interface/lambdaapigw"
	"app/internal/driver_interface_adapter/lambdaapigw_adapter"
	"sync"
)

func NewDriverInterfaceMiddlewareRegistry(
	useCaseRegistry *UseCaseRegistry,
) *DriverInterfaceMiddlewareRegistry {
	registry := DriverInterfaceMiddlewareRegistry{}
	registry.Initialize(useCaseRegistry)
	return &registry
}

type DriverInterfaceMiddlewareRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry

	authenticateBeforeMiddleware lambdaapigw.BeforeMiddleware
}

func (registry *DriverInterfaceMiddlewareRegistry) Initialize(useCaseRegistry *UseCaseRegistry) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry
	registry.authenticateBeforeMiddleware = lambdaapigw_adapter.NewAuthenticateBeforeMiddleware(
		registry.useCaseRegistry.BankAccountAuthMWUseCase(),
	)
	registry.Unlock()
}

func (registry *DriverInterfaceMiddlewareRegistry) AuthenticateBeforeMiddleware() lambdaapigw.BeforeMiddleware {
	if registry.authenticateBeforeMiddleware == nil {
		registry.Lock()
		registry.authenticateBeforeMiddleware = lambdaapigw_adapter.NewAuthenticateBeforeMiddleware(
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
