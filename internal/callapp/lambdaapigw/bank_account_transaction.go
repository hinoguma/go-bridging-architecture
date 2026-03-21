package lambdaapigw

import (
	"fmt"
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
	Amount              *int64  `json:"amount,omitempty"`
	Currency            *string `json:"currency,omitempty"`
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
	Amount              *int64  `json:"amount,omitempty"`
	Currency            *string `json:"currency,omitempty"`
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
