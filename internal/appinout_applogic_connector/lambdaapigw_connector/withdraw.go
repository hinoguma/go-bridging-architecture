package lambdaapigw_connector

import (
	"app/internal/appinout/lambdaapigw"
	"app/internal/applogic/domain/model"
	"app/internal/applogic/usecase"
	"app/internal/crosscutting/errors"
	"context"
)

func NewWithdrawHandlerConnector(
	uc usecase.WithdrawUseCase,
) lambdaapigw.LambdaAPIGWHandler {
	return WithdrawHandlerConnector{uc: uc}
}

type WithdrawHandlerConnector struct {
	uc usecase.WithdrawUseCase
}

func (handler WithdrawHandlerConnector) Do(ctx context.Context, request lambdaapigw.HandlerRequest) (lambdaapigw.HandlerResponse, error) {
	internalServerErr := lambdaapigw.NewInternalServerErrorResponse()

	body, err := lambdaapigw.NewWithdrawRequestBody(request)
	if err != nil {
		return lambdaapigw.NewBodyDecodingErrorResponse(), errors.LiftWithCtx(err, ctx)
	}

	validateDetails := lambdaapigw.ValidateWithdrawRequestBody(
		body, model.MaxWithdrawAmountOneTime, []string{model.JPY.String()},
	)
	if len(validateDetails) > 0 {
		return lambdaapigw.NewValidateErrorResponseFromValidateDetails(validateDetails), nil
	}

	bodyDTO := WithdrawRequestBodyDTO(body)

	input := ConvertToWithdrawUseCaseInput(request, bodyDTO)

	ucOutput, err := handler.uc.Do(ctx, input)
	if err != nil {
		return internalServerErr, errors.Lift(err)
	}

	resp := ConvertWithdrawUseCaseOutputToHandlerResponse(ucOutput)
	return resp, nil
}

func ConvertToWithdrawUseCaseInput(req lambdaapigw.HandlerRequest, body WithdrawRequestBodyDTO) usecase.WithdrawUseCaseInput {
	return usecase.WithdrawUseCaseInput{
		BankAccountID: GetBankAccountIDFromAuthenticatedRequest(req),
		TransactionID: body.GetTransactionRecordID(),
		Amount:        body.Money(),
	}
}

func ConvertWithdrawUseCaseOutputToHandlerResponse(output usecase.WithdrawUseCaseOutput) lambdaapigw.HandlerResponse {
	if output.NotEnoughBalance {
		return lambdaapigw.NewNotEnoughBalanceResponse()
	}
	return lambdaapigw.NewWithdrawSuccessResponse(
		output.TransactionRecord.ID.String(),
	)
}

type WithdrawRequestBodyDTO lambdaapigw.WithdrawRequestBody

func (body WithdrawRequestBodyDTO) GetAmount() float64 {
	if body.Amount == nil {
		return 0
	}
	return float64(*body.Amount)
}

func (body WithdrawRequestBodyDTO) GetCurrency() model.Currency {
	if body.Currency == nil {
		return ""
	}
	return model.Currency(*body.Currency)
}

func (body WithdrawRequestBodyDTO) Money() model.Money {
	return model.Money{
		Amount:   body.GetAmount(),
		Currency: body.GetCurrency(),
	}
}

func (body WithdrawRequestBodyDTO) GetTransactionRecordID() *model.TransactionRecordID {
	if body.TransactionRecordID == nil {
		return nil
	}
	id := model.TransactionRecordID(*body.TransactionRecordID)
	return &id
}
