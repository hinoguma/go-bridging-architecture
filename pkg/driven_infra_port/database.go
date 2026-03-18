package driven_infra_port

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/driven_infra/rdb"
)

func convertModelToSQLOperationOptions(
	modelOptions model.DBOperationOptions,
) rdb.SQLOperationOptions {
	options := rdb.SQLOperationOptions{}
	if modelOptions.HasDBTransactionID() {
		options.SetTransactionID(
			modelOptions.GetTransactionID().String(),
		)
	}
	return options
}

func convertOptionalFuncsToSQLOperationOptions(optionalFuncs []model.DBOperationOptionalFunc) rdb.SQLOperationOptions {
	return convertModelToSQLOperationOptions(
		model.ApplyDBOperationOptionalFuncs(optionalFuncs...),
	)
}
