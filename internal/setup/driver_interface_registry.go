package setup

import (
	"app/internal/callapp/lambdaapigw"
	"app/internal/callapp_applogic_bridge/lambdaapigw_bridge"
	"context"
	"sync"
)

type DriverInterfaceRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry
	depositHandler  lambdaapigw.LambdaAPIGWBHandler
	withdrawHandler lambdaapigw.LambdaAPIGWBHandler
}

func (registry *DriverInterfaceRegistry) Initialize(
	ctx context.Context,
	useCaseRegistry *UseCaseRegistry,
) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry

	handler := lambdaapigw_bridge.NewDepositHandlerLambdaAPIGateway(
		registry.useCaseRegistry.DepositUseCase(),
	)
	registry.depositHandler = handler

	registry.Unlock()
}

func (registry *DriverInterfaceRegistry) DepositHandler() lambdaapigw.LambdaAPIGWBHandler {
	if registry.depositHandler == nil {
		registry.Lock()
		handler := lambdaapigw_bridge.NewDepositHandlerLambdaAPIGateway(
			registry.useCaseRegistry.DepositUseCase(),
		)
		registry.depositHandler = handler
		registry.Unlock()
	}
	return registry.depositHandler
}

func (registry *DriverInterfaceRegistry) WithdrawHandler() lambdaapigw.LambdaAPIGWBHandler {
	if registry.withdrawHandler == nil {
		registry.Lock()
		handler := lambdaapigw_bridge.NewWithdrawHandlerLambdaAPIGateway(
			registry.useCaseRegistry.WithdrawUseCase(),
		)
		registry.withdrawHandler = handler
		registry.Unlock()
	}
	return registry.withdrawHandler
}

func GetDriverInterfaceRegistry() *DriverInterfaceRegistry {
	return &driverInterfaceRegistry
}

var driverInterfaceRegistry DriverInterfaceRegistry = DriverInterfaceRegistry{}
