package setup

import (
	"app/internal/callapp/lambdaapigw"
	"app/internal/callapp_applogic_connector/lambdaapigw_connector"
	"context"
	"sync"
)

type CallAppRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry
	depositHandler  lambdaapigw.LambdaAPIGWHandler
	withdrawHandler lambdaapigw.LambdaAPIGWHandler
}

func (registry *CallAppRegistry) Initialize(
	ctx context.Context,
	useCaseRegistry *UseCaseRegistry,
) {
	registry.Lock()
	registry.useCaseRegistry = useCaseRegistry

	handler := lambdaapigw_connector.NewDepositHandlerConnector(
		registry.useCaseRegistry.DepositUseCase(),
	)
	registry.depositHandler = handler

	registry.Unlock()
}

func (registry *CallAppRegistry) DepositHandler() lambdaapigw.LambdaAPIGWHandler {
	if registry.depositHandler == nil {
		registry.Lock()
		handler := lambdaapigw_connector.NewDepositHandlerConnector(
			registry.useCaseRegistry.DepositUseCase(),
		)
		registry.depositHandler = handler
		registry.Unlock()
	}
	return registry.depositHandler
}

func (registry *CallAppRegistry) WithdrawHandler() lambdaapigw.LambdaAPIGWHandler {
	if registry.withdrawHandler == nil {
		registry.Lock()
		handler := lambdaapigw_connector.NewWithdrawHandlerConnector(
			registry.useCaseRegistry.WithdrawUseCase(),
		)
		registry.withdrawHandler = handler
		registry.Unlock()
	}
	return registry.withdrawHandler
}

func GetCallAppRegistry() *CallAppRegistry {
	return &callAppRegistry
}

var callAppRegistry CallAppRegistry = CallAppRegistry{}
