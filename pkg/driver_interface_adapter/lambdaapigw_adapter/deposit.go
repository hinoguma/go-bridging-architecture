package lambdaapigw_adapter

import (
	"app/pkg/applogic/domain/model"
	"app/pkg/applogic/usecase"
	"app/pkg/crosscutting"
	"app/pkg/crosscutting/errors"
	"app/pkg/driver_interface/lambdaapigw"
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

	body, err := DecodeEventToBody(request)
	if err != nil {
		return lambdaapigw.NewBodyDecodingErrorResponse(), errors.Lift(err)
	}

	validateDetails := ValidateDepositRequestBody(body)
	if len(validateDetails) > 0 {
		return lambdaapigw.NewValidateErrorResponseFromValidateDetails(validateDetails), nil
	}

	input := ConvertBodyToUseCaseInput(body)

	ucOutput, err := handler.uc.Do(ctx, input)
	if err != nil {
		return internalServerErr, errors.Lift(err)
	}

	resp := ConvertDepositUseCaseOutputToHandlerResponse(ucOutput)
	return resp, nil
}

func DecodeEventToBody(request lambdaapigw.HandlerRequest) (DepositRequestBody, error) {
	body := DepositRequestBody{}
	err := json.Unmarshal([]byte(request.Raw.Body), &body)
	if err != nil {
		return body, errors.Lift(err)
	}
	return body, nil
}

func ConvertBodyToUseCaseInput(body DepositRequestBody) usecase.DepositUseCaseInput {
	return usecase.DepositUseCaseInput{
		BankAccountID: "",
		Amount:        body.Money(),
	}
}

func ConvertDepositUseCaseOutputToHandlerResponse(output usecase.DepositUseCaseOutput) lambdaapigw.HandlerResponse {
	return lambdaapigw.NewSuccessResponse(
		fmt.Sprintf("{\"transactionId\": \"%s\"}", output.TransactionRecord.ID.String()),
	)
}

type DepositRequestBody struct {
	Amount   *int64  `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
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
	}

	if body.Currency == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "currency",
		})
	} else if !body.GetCurrency().IsValid() {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeFormat,
			Field: "currency",
		})
	}
	return details
}
