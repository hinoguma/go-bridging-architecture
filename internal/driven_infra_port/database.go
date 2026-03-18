package driven_infra_port

import (
	"app/internal/applogic/domain/model"
	"app/internal/driven_infra/postgres"
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
