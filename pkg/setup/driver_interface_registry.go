package setup

import (
	"app/pkg/driver_interface/lambdaapigw"
	"app/pkg/driver_interface_adapter/lambdaapigw_adapter"
	"context"
	"sync"
)

type DriverInterfaceRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry
	depositHandler  lambdaapigw.LambdaAPIGWBHandler
}

func (registry *DriverInterfaceRegistry) Initialize(
	ctx context.Context,
	useCaseRegistry *UseCaseRegistry,
) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry

	handler := lambdaapigw_adapter.NewDepositHandlerLambdaAPIGateway(
		registry.useCaseRegistry.DepositUseCase(),
	)
	registry.depositHandler = handler

	registry.Unlock()
}

func (registry *DriverInterfaceRegistry) DepositHandler() lambdaapigw.LambdaAPIGWBHandler {
	if registry.depositHandler == nil {
		registry.Lock()
		handler := lambdaapigw_adapter.NewDepositHandlerLambdaAPIGateway(
			registry.useCaseRegistry.DepositUseCase(),
		)
		registry.depositHandler = handler
		registry.Unlock()
	}
	return registry.depositHandler
}

func GetDriverInterfaceRegistry() *DriverInterfaceRegistry {
	return &driverInterfaceRegistry
}

var driverInterfaceRegistry DriverInterfaceRegistry = DriverInterfaceRegistry{}
