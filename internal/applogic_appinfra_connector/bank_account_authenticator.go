package applogic_appinfra_connector

import (
	"app/internal/appinfra/aws"
	"app/internal/applogic/domain/model"
	"app/internal/applogic/domain/repository"
	"app/internal/crosscutting/errors"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type bankAccountAuthenticator struct {
	cognitoClient driven_infra_aws.CognitoClient
}

func NewBankAccountAuthenticator(
	cognitoClient driven_infra_aws.CognitoClient,
) repository.BankAccountAuthenticator {
	return bankAccountAuthenticator{
		cognitoClient: cognitoClient,
	}
}

func (authenticator bankAccountAuthenticator) VerifyToken(
	ctx context.Context, req model.BankAccountVerifyTokenRequest,
) model.BankAccountVerifyTokenResult {
	result := model.BankAccountVerifyTokenResult{}
	res := authenticator.cognitoClient.VerifyAccessToken(
		ctx, "", driven_infra_aws.CognitoAccessTokenJWT(req.Token), req.GetRequestAt(),
	)
	if res.IsErr() {
		result.Err = errors.LiftWithCtx(res.Err, ctx)
		return result
	}
	if !res.IsValid() {
		if !res.IsValidExp {
			result.IsTokenExpired = true
		} else {
			result.IsInValidToken = true
		}
		return result
	}
	userInfo, err := authenticator.cognitoClient.GetUserInfo(ctx, req.Token)
	if err != nil {
		result.Err = errors.LiftWithCtx(err, ctx)
		return result
	}
	result.BankAccountID, err = GetBankAccountIDFromUserInfo(userInfo)
	if err != nil {
		result.Err = errors.LiftWithCtx(err, ctx)
		return result
	}
	return result
}

func GetBankAccountIDFromUserInfo(
	userInfo cognitoidentityprovider.GetUserOutput,
) (model.BankAccountID, error) {

	for _, attr := range userInfo.UserAttributes {
		if attr.Name == nil || attr.Value == nil {
			continue
		}
		if *attr.Name == "custom:bank_account_id" {
			return model.BankAccountID(*attr.Value), nil
		}
	}
	return "", errors.New("bank_account_id not found in user attributes")
}
