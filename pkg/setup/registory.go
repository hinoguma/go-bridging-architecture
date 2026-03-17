package setup

import (
	"app/pkg/applogic/domain/repository"
	"app/pkg/applogic/usecase"
	"app/pkg/driven_infra/aws"
	"app/pkg/driven_infra_port"
	"app/pkg/driver_interface/lambdaapigw"
	"app/pkg/driver_interface_adapter/lambdaapigw_adapter"
	"context"
)

type DrivenInfraRegistry struct {
	DynamoDBClient driven_infra_aws.DynamoDBClient
}

func (registry DrivenInfraRegistry) GetDynamoDBClient() driven_infra_aws.DynamoDBClient {
	if registry.DynamoDBClient == nil {

	}
	return registry.DynamoDBClient
}

type DomainRepositoryRegistry struct {
	DrivenInfraRegistry         DrivenInfraRegistry
	BankAccountRepository       repository.BankAccountRepository
	TransactionRecordRepository repository.TransactionRecordRepository
}

func (registry DomainRepositoryRegistry) GetBankAccountRepository() repository.BankAccountRepository {
	if registry.BankAccountRepository == nil {
		registry.BankAccountRepository = driven_infra_port.NewBankAccountRepositoryDynamoDB(
			registry.DrivenInfraRegistry.GetDynamoDBClient(),
		)
	}
	return registry.BankAccountRepository
}

func (registry DomainRepositoryRegistry) GetTransactionRecordRepository() repository.TransactionRecordRepository {
	if registry.TransactionRecordRepository == nil {
		registry.TransactionRecordRepository = driven_infra_port.NewTransactionRecordRepositoryDynamoDB(
			registry.DrivenInfraRegistry.GetDynamoDBClient(),
		)
	}
	return registry.TransactionRecordRepository
}

type UseCaseRegistry struct {
	DomainRepositoryRegistry DomainRepositoryRegistry

	DepositUseCase usecase.DepositUseCase
}

func (registry UseCaseRegistry) GetDepositUseCase() usecase.DepositUseCase {
	if registry.DepositUseCase == nil {
		registry.DepositUseCase = usecase.NewDepositUseCase(
			registry.DomainRepositoryRegistry.GetBankAccountRepository(),
			registry.DomainRepositoryRegistry.GetTransactionRecordRepository(),
		)
	}
	return registry.DepositUseCase
}

type DriverInterfaceRegistry struct {
	useCaseRegistry UseCaseRegistry
	depositHandler  lambdaapigw.LambdaAPIGWBHandler
}

func (d *DriverInterfaceRegistry) Initialize(
	ctx context.Context,
	useCaseRegistry UseCaseRegistry,
) {
	d.useCaseRegistry = useCaseRegistry
}

func (d *DriverInterfaceRegistry) GetDepositHandler() lambdaapigw.LambdaAPIGWBHandler {
	if d.depositHandler == nil {
		handler := lambdaapigw_adapter.NewDepositHandlerLambdaAPIGateway(
			d.useCaseRegistry.GetDepositUseCase(),
		)
		d.depositHandler = handler
	}
	return d.depositHandler
}

func GetDriverInterfaceRegistry() *DriverInterfaceRegistry {
	return &driverInterfaceRegistry
}

var driverInterfaceRegistry DriverInterfaceRegistry = DriverInterfaceRegistry{}
