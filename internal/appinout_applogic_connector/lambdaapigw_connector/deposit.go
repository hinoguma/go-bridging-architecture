package lambdaapigw_connector

import (
	"app/internal/appinout/lambdaapigw"
	"app/internal/applogic/domain/model"
	"app/internal/applogic/usecase"
	"app/internal/crosscutting/errors"
	"context"
)

func NewDepositHandlerConnector(
	uc usecase.DepositUseCase,
) lambdaapigw.LambdaAPIGWHandler {
	return DepositHandlerConnector{uc: uc}
}

type DepositHandlerConnector struct {
	uc usecase.DepositUseCase
}

func (handler DepositHandlerConnector) Do(ctx context.Context, request lambdaapigw.HandlerRequest) (lambdaapigw.HandlerResponse, error) {
	internalServerErr := lambdaapigw.NewInternalServerErrorResponse()

	body, err := lambdaapigw.NewDepositRequestBody(request)
	if err != nil {
		return lambdaapigw.NewBodyDecodingErrorResponse(), errors.LiftWithCtx(err, ctx)
	}

	validateDetails := lambdaapigw.ValidateDepositRequestBody(
		body, model.MaxDepositAmountOneTime, []string{model.JPY.String()},
	)
	if len(validateDetails) > 0 {
		return lambdaapigw.NewValidateErrorResponseFromValidateDetails(validateDetails), nil
	}

	bodyDTO := DepositRequestBodyDTO(body)

	input := ConvertToDepositUseCaseInput(request, bodyDTO)

	ucOutput, err := handler.uc.Do(ctx, input)
	if err != nil {
		return internalServerErr, errors.Lift(err)
	}

	resp := ConvertDepositUseCaseOutputToHandlerResponse(ucOutput)
	return resp, nil
}

func ConvertToDepositUseCaseInput(req lambdaapigw.HandlerRequest, body DepositRequestBodyDTO) usecase.DepositUseCaseInput {
	return usecase.DepositUseCaseInput{
		BankAccountID: GetBankAccountIDFromAuthenticatedRequest(req),
		TransactionID: body.GetTransactionRecordID(),
		Amount:        body.Money(),
	}
}

func ConvertDepositUseCaseOutputToHandlerResponse(output usecase.DepositUseCaseOutput) lambdaapigw.HandlerResponse {
	return lambdaapigw.NewDepositSuccessResponse(
		output.TransactionRecord.ID.String(),
	)
}

type DepositRequestBodyDTO lambdaapigw.DepositRequestBody

func (body DepositRequestBodyDTO) GetAmount() float64 {
	if body.Amount == nil {
		return 0
	}
	return float64(*body.Amount)
}

func (body DepositRequestBodyDTO) GetCurrency() model.Currency {
	if body.Currency == nil {
		return ""
	}
	return model.Currency(*body.Currency)
}

func (body DepositRequestBodyDTO) Money() model.Money {
	return model.Money{
		Amount:   body.GetAmount(),
		Currency: body.GetCurrency(),
	}
}

func (body DepositRequestBodyDTO) GetTransactionRecordID() *model.TransactionRecordID {
	if body.TransactionRecordID == nil {
		return nil
	}
	id := model.TransactionRecordID(*body.TransactionRecordID)
	return &id
}
