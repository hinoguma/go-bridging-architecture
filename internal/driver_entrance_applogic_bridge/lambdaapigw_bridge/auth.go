package lambdaapigw_bridge

import (
	"app/internal/applogic/domain/model"
	"app/internal/applogic/usecase"
	"app/internal/crosscutting/errors"
	"app/internal/driver_entrance/lambdaapigw"
	"context"
)

func NewAuthenticateBeforeMiddleware(
	uc usecase.BankAccountAuthMWUseCase,
) lambdaapigw.BeforeMiddleware {
	return authenticateBeforeMiddleware{
		uc: uc,
	}
}

type authenticateBeforeMiddleware struct {
	uc usecase.BankAccountAuthMWUseCase
}

func (mw authenticateBeforeMiddleware) Do(
	ctx context.Context,
	request lambdaapigw.HandlerRequest,
) (lambdaapigw.BeforeMiddlewareResult, error) {

	result := lambdaapigw.BeforeMiddlewareResult{
		Ctx:      ctx,
		Request:  request,
		Response: nil,
	}
	mwReq := ConvertHandlerRequestToBankAccountAuthMWRequest(request)
	mwRes := mw.uc.Do(ctx, mwReq)
	if mwRes.IsErr() {
		return result, errors.Lift(mwRes.Err)
	}

	if mwRes.IsTokenExpired {
		result.SetResponse(
			lambdaapigw.NewAuthTokenExpiredErrorResponse(),
		)
		return result, nil
	}
	if mwRes.IsInValidToken {
		result.SetResponse(
			lambdaapigw.NewAuthTokenExpiredErrorResponse(),
		)
		return result, nil
	}

	result.Request.SetAuthenticatedBankAccountID(
		mwRes.BankAccount.ID.String(),
	)

	return result, nil
}

func ConvertHandlerRequestToBankAccountAuthMWRequest(request lambdaapigw.HandlerRequest) usecase.BankAccountAuthMWRequest {
	v, _ := request.GetAuthorization()
	return usecase.BankAccountAuthMWRequest{
		Token: v,
	}
}

func GetBankAccountIDFromAuthenticatedRequest(request lambdaapigw.HandlerRequest) model.BankAccountID {
	return model.BankAccountID(request.GetAuthenticatedBankAccountID())
}
