package lambdaapigw_adapter

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/usecase"
	"app/internal/crosscutting"
	"app/internal/crosscutting/errors"
	"app/internal/driver_interface/lambdaapigw"
	"context"
	"encoding/json"
	"fmt"
)

func NewDepositHandlerLambdaAPIGateway(
	uc usecase.DepositUseCase,
) lambdaapigw.LambdaAPIGWBHandler {
	return DepositHandlerLambdaAPIGateway{uc: uc}
}

type DepositHandlerLambdaAPIGateway struct {
	uc usecase.DepositUseCase
}

func (handler DepositHandlerLambdaAPIGateway) Do(ctx context.Context, request lambdaapigw.HandlerRequest) (lambdaapigw.HandlerResponse, error) {
	internalServerErr := lambdaapigw.NewInternalServerErrorResponse()

	body, err := DecodeEventToDepositRequestBody(request)
	if err != nil {
		return lambdaapigw.NewBodyDecodingErrorResponse(), errors.LiftWithCtx(err, ctx)
	}

	validateDetails := ValidateDepositRequestBody(body)
	if len(validateDetails) > 0 {
		return lambdaapigw.NewValidateErrorResponseFromValidateDetails(validateDetails), nil
	}

	input := ConvertToDepositUseCaseInput(request, body)

	ucOutput, err := handler.uc.Do(ctx, input)
	if err != nil {
		return internalServerErr, errors.Lift(err)
	}

	resp := ConvertDepositUseCaseOutputToHandlerResponse(ucOutput)
	return resp, nil
}

func DecodeEventToDepositRequestBody(request lambdaapigw.HandlerRequest) (DepositRequestBody, error) {
	body := DepositRequestBody{}
	err := json.Unmarshal([]byte(request.Raw.Body), &body)
	if err != nil {
		return body, errors.Lift(err)
	}
	return body, nil
}

func ConvertToDepositUseCaseInput(req lambdaapigw.HandlerRequest, body DepositRequestBody) usecase.DepositUseCaseInput {
	return usecase.DepositUseCaseInput{
		BankAccountID: GetBankAccountIDFromAuthenticatedRequest(req),
		TransactionID: body.GetTransactionRecordID(),
		Amount:        body.Money(),
	}
}

func ConvertDepositUseCaseOutputToHandlerResponse(output usecase.DepositUseCaseOutput) lambdaapigw.HandlerResponse {
	return lambdaapigw.NewSuccessResponse(
		fmt.Sprintf("{\"transactionId\": \"%s\"}", output.TransactionRecord.ID.String()),
	)
}

type DepositRequestBody struct {
	TransactionRecordID *string `json:"transactionRecordId,omitempty"`
	Amount              *int64  `json:"amount,omitempty"`
	Currency            *string `json:"currency,omitempty"`
}

func (body DepositRequestBody) GetAmount() int64 {
	if body.Amount == nil {
		return 0
	}
	return *body.Amount
}

func (body DepositRequestBody) GetCurrency() model.Currency {
	if body.Currency == nil {
		return ""
	}
	return model.Currency(*body.Currency)
}

func (body DepositRequestBody) Money() model.Money {
	return model.Money{
		Amount:   body.GetAmount(),
		Currency: body.GetCurrency(),
	}
}

func (body DepositRequestBody) GetTransactionRecordID() *model.TransactionRecordID {
	if body.TransactionRecordID == nil {
		return nil
	}
	id := model.TransactionRecordID(*body.TransactionRecordID)
	return &id
}

func ValidateDepositRequestBody(body DepositRequestBody) []crosscutting.ValidateDetail {
	details := make([]crosscutting.ValidateDetail, 0)

	if body.Amount == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "amount",
		})
	} else if body.GetAmount() <= 0 {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeMin,
			Field: "amount",
			Min:   crosscutting.Ptr(0),
		})
	} else if body.GetAmount() > model.MaxWithdrawAmountOneTime {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeMax,
			Field: "amount",
			Max:   crosscutting.Ptr(model.MaxWithdrawAmountOneTime),
		})
	}

	if body.Currency == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "currency",
		})
	} else if body.GetCurrency() != model.JPY {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeFormat,
			Field: "currency",
		})
	}

	if body.TransactionRecordID == nil {
		// transactionRecordId is optional

	} else if *body.TransactionRecordID == "" {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeMinStrLen,
			Field: "transactionRecordId",
			Min:   crosscutting.Ptr(1),
		})
	}
	return details
}
