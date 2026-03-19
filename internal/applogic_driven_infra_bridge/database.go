package applogic_driven_infra_bridge

import (
	"app/internal/applogic/domain/model"
	"app/internal/appinfra/postgres"
)

func convertModelToSQLOperationOptions(
	modelOptions model.DBOperationOptions,
) postgres.SQLOperationOptions {
	options := postgres.SQLOperationOptions{}
	if modelOptions.HasDBTransactionID() {
		options.SetTransactionID(
			modelOptions.GetTransactionID().String(),
		)
	}
	if modelOptions.HasSelectLock() && modelOptions.GetSelectLock() {
		options.SetLockMode(postgres.LockForUpdate)
	}
	return options
}

func convertOptionalFuncsToSQLOperationOptions(optionalFuncs []model.DBOperationOptionalFunc) postgres.SQLOperationOptions {
	return convertModelToSQLOperationOptions(
		model.ApplyDBOperationOptionalFuncs(optionalFuncs...),
	)
}
