package lambdaapigw_bridge

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/usecase"
	"app/internal/callapp/lambdaapigw"
	"app/internal/crosscutting"
	"app/internal/crosscutting/errors"
	"context"
	"encoding/json"
)

func NewWithdrawHandlerLambdaAPIGateway(
	uc usecase.WithdrawUseCase,
) lambdaapigw.LambdaAPIGWBHandler {
	return WithdrawHandlerLambdaAPIGateway{uc: uc}
}

type WithdrawHandlerLambdaAPIGateway struct {
	uc usecase.WithdrawUseCase
}

func (handler WithdrawHandlerLambdaAPIGateway) Do(ctx context.Context, request lambdaapigw.HandlerRequest) (lambdaapigw.HandlerResponse, error) {
	internalServerErr := lambdaapigw.NewInternalServerErrorResponse()

	body, err := DecodeEventToWithdrawRequestBodyDTO(request)
	if err != nil {
		return lambdaapigw.NewBodyDecodingErrorResponse(), errors.LiftWithCtx(err, ctx)
	}

	validateDetails := ValidateWithdrawRequestBodyDTO(body)
	if len(validateDetails) > 0 {
		return lambdaapigw.NewValidateErrorResponseFromValidateDetails(validateDetails), nil
	}

	input := ConvertToWithdrawUseCaseInput(request, body)

	ucOutput, err := handler.uc.Do(ctx, input)
	if err != nil {
		return internalServerErr, errors.Lift(err)
	}

	resp := ConvertWithdrawUseCaseOutputToHandlerResponse(ucOutput)
	return resp, nil
}

func DecodeEventToWithdrawRequestBodyDTO(request lambdaapigw.HandlerRequest) (WithdrawRequestBodyDTO, error) {
	body := WithdrawRequestBodyDTO{}
	err := json.Unmarshal([]byte(request.Raw.Body), &body)
	if err != nil {
		return body, errors.Lift(err)
	}
	return body, nil
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

func (body WithdrawRequestBodyDTO) GetAmount() int64 {
	if body.Amount == nil {
		return 0
	}
	return *body.Amount
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

func ValidateWithdrawRequestBodyDTO(body WithdrawRequestBodyDTO) []crosscutting.ValidateDetail {
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
