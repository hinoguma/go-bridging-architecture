package lambdaapigw

import (
	"app/internal/crosscutting"
	"app/internal/crosscutting/errors"
	"encoding/json"
	"fmt"
	"slices"
)

/******************************
 * Shared
 ******************************/

type TransactionSuccessResponseBody struct {
	TransactionID string `json:"transactionId"`
}

func (body TransactionSuccessResponseBody) JsonString() string {
	return fmt.Sprintf(`{"transactionId":"%s"}`, body.TransactionID)
}

func NewTransactionSuccessResponse(transactionID string) HandlerResponse {
	return NewSuccessResponse(
		TransactionSuccessResponseBody{TransactionID: transactionID}.JsonString(),
	)
}

/******************************
 * POST /bank/account/deposit
 ******************************/

type DepositRequestBody struct {
	TransactionRecordID *string `json:"transactionRecordId,omitempty"`
	Amount              *int    `json:"amount,omitempty"`
	Currency            *string `json:"currency,omitempty"`
}

func NewDepositRequestBody(request HandlerRequest) (DepositRequestBody, error) {
	body := DepositRequestBody{}
	if request.Raw.Body == "" {
		return body, nil
	}
	err := json.Unmarshal([]byte(request.Raw.Body), &body)
	if err != nil {
		return body, errors.Lift(err)
	}
	return body, nil
}

func ValidateDepositRequestBody(
	body DepositRequestBody,
	maxAmountOneTime int,
	allowedCurrencies []string,
) []crosscutting.ValidateDetail {
	details := make([]crosscutting.ValidateDetail, 0)

	if body.Amount == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "amount",
		})
	} else if *body.Amount <= 0 {
		details = append(details, crosscutting.NewMinValidateDetail(
			"amount", 0,
		))
	} else if *body.Amount > maxAmountOneTime {
		details = append(details, crosscutting.NewMaxValidateDetail(
			"amount", float64(maxAmountOneTime),
		))
	}

	if body.Currency == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "currency",
		})
	} else if slices.Contains(allowedCurrencies, *body.Currency) == false {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeFormat,
			Field: "currency",
		})
	}

	if body.TransactionRecordID == nil {
		// transactionRecordId is optional

	} else if *body.TransactionRecordID == "" {
		details = append(details, crosscutting.NewMinStrLenValidateDetail(
			"transactionRecordId", 1,
		))
	}
	return details
}

type DepositSuccessResponseBody TransactionSuccessResponseBody

func NewDepositSuccessResponse(transactionID string) HandlerResponse {
	return NewTransactionSuccessResponse(transactionID)
}

/******************************
 * POST /bank/account/withdraw
 ******************************/

const NotEnoughBalanceAPIStatus int = 460

type WithdrawRequestBody struct {
	TransactionRecordID *string `json:"transactionRecordId,omitempty"`
	Amount              *int    `json:"amount,omitempty"`
	Currency            *string `json:"currency,omitempty"`
}

func NewWithdrawRequestBody(request HandlerRequest) (WithdrawRequestBody, error) {
	body := WithdrawRequestBody{}
	if request.Raw.Body == "" {
		return body, nil
	}
	err := json.Unmarshal([]byte(request.Raw.Body), &body)
	if err != nil {
		return body, errors.Lift(err)
	}
	return body, nil
}

func ValidateWithdrawRequestBody(
	body WithdrawRequestBody,
	maxAmountOneTime int,
	allowedCurrencies []string,
) []crosscutting.ValidateDetail {
	details := make([]crosscutting.ValidateDetail, 0)

	if body.Amount == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "amount",
		})
	} else if *body.Amount <= 0 {
		details = append(details, crosscutting.NewMinValidateDetail(
			"amount", 0,
		))
	} else if *body.Amount > maxAmountOneTime {
		details = append(details, crosscutting.NewMaxValidateDetail(
			"amount", float64(maxAmountOneTime),
		))
	}

	if body.Currency == nil {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeRequired,
			Field: "currency",
		})
	} else if slices.Contains(allowedCurrencies, *body.Currency) == false {
		details = append(details, crosscutting.ValidateDetail{
			Type:  crosscutting.InValidTypeFormat,
			Field: "currency",
		})
	}

	if body.TransactionRecordID == nil {
		// transactionRecordId is optional

	} else if *body.TransactionRecordID == "" {
		details = append(details, crosscutting.NewMinStrLenValidateDetail(
			"transactionRecordId", 1,
		))
	}
	return details
}

func NewNotEnoughBalanceResponse() HandlerResponse {
	return NewErrorResponse(
		NotEnoughBalanceAPIStatus,
		"your balance is not enough for this withdraw",
		nil,
	)
}

type WithdrawSuccessResponseBody TransactionSuccessResponseBody

func NewWithdrawSuccessResponse(transactionID string) HandlerResponse {
	return NewTransactionSuccessResponse(transactionID)
}
