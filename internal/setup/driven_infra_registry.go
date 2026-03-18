package setup

import (
	"app/internal/crosscutting/errors"
	"app/internal/driven_infra/aws"
	"app/internal/driven_infra/postgres"
	"context"
	"database/sql"
	"sync"
)

type DrivenInfraRegistry struct {
	sync.Mutex
	dynamoDBClient     driven_infra_aws.DynamoDBClient
	cognitoClient      driven_infra_aws.CognitoClient
	sqlClient          postgres.SQLClient
	transactionManager postgres.TransactionManagerIF
}

func (registry *DrivenInfraRegistry) Initialize(
	ctx context.Context,
	db *sql.DB,
) error {
	var err error
	registry.Lock()
	registry.dynamoDBClient, err = driven_infra_aws.NewDynamoDBClient(ctx)
	if err != nil {
		registry.Unlock()
		return errors.Lift(err)
	}
	registry.cognitoClient, err = driven_infra_aws.NewCognitoClient(ctx)
	if err != nil {
		registry.Unlock()
		return errors.Lift(err)
	}
	registry.sqlClient = postgres.NewPostgreSQLClient(db)
	registry.transactionManager = postgres.NewTransactionManager(registry.sqlClient.GetDB())
	registry.Unlock()
	return nil
}

func (registry *DrivenInfraRegistry) DynamoDBClient() driven_infra_aws.DynamoDBClient {
	return registry.dynamoDBClient
}

func (registry *DrivenInfraRegistry) CognitoClient() driven_infra_aws.CognitoClient {
	return registry.cognitoClient
}

func (registry *DrivenInfraRegistry) SQLClient() postgres.SQLClient {
	return registry.sqlClient
}

func (registry *DrivenInfraRegistry) TransactionManager() postgres.TransactionManagerIF {
	return registry.transactionManager
}

var drivenInfraRegistry = DrivenInfraRegistry{}

func GetDrivenInfraRegistry() *DrivenInfraRegistry {
	return &drivenInfraRegistry
}
