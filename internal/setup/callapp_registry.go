package setup

import (
	"app/internal/appinout/lambdaapigw"
	"app/internal/appinout_applogic_connector/lambdaapigw_connector"
	"context"
	"sync"
)

type AppInOutRegistry struct {
	sync.Mutex
	useCaseRegistry *UseCaseRegistry
	depositHandler  lambdaapigw.LambdaAPIGWHandler
	withdrawHandler lambdaapigw.LambdaAPIGWHandler
}

func (registry *AppInOutRegistry) Initialize(
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

func (registry *AppInOutRegistry) DepositHandler() lambdaapigw.LambdaAPIGWHandler {
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

func (registry *AppInOutRegistry) WithdrawHandler() lambdaapigw.LambdaAPIGWHandler {
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

func GetAppInOutRegistry() *AppInOutRegistry {
	return &appInOutRegistry
}

var appInOutRegistry AppInOutRegistry = AppInOutRegistry{}
