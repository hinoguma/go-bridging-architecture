package setup

import (
	"app/pkg/crosscutting/errors"
	"app/pkg/driven_infra/aws"
	"app/pkg/driven_infra/rdb"
	"context"
	"database/sql"
	"sync"
)

type DrivenInfraRegistry struct {
	sync.Mutex
	dynamoDBClient     driven_infra_aws.DynamoDBClient
	cognitoClient      driven_infra_aws.CognitoClient
	sqlClient          rdb.SQLClient
	transactionManager rdb.TransactionManagerIF
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
	registry.sqlClient = rdb.NewPostgreSQLClient(db)
	rdb.NewTransactionManager(registry.sqlClient.GetDB())
	registry.Unlock()
	return nil
}

func (registry *DrivenInfraRegistry) DynamoDBClient() driven_infra_aws.DynamoDBClient {
	return registry.dynamoDBClient
}

func (registry *DrivenInfraRegistry) CognitoClient() driven_infra_aws.CognitoClient {
	return registry.cognitoClient
}

func (registry *DrivenInfraRegistry) SQLClient() rdb.SQLClient {
	return registry.sqlClient
}

func (registry *DrivenInfraRegistry) TransactionManager() rdb.TransactionManagerIF {
	return registry.transactionManager
}

var drivenInfraRegistry = DrivenInfraRegistry{}

func GetDrivenInfraRegistry() *DrivenInfraRegistry {
	return &drivenInfraRegistry
}
