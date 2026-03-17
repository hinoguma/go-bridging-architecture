package driven_infra_aws

import (
	"app/pkg/crosscutting/errors"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoClient interface {
	VerifyAccessToken(
		ctx context.Context,
		userPoolID string,
		token CognitoAccessTokenJWT,
		requestAt time.Time,
	) VerifyTokenResult[CognitoAccessTokenJWTPayload]

	GetUserInfo(
		ctx context.Context,
		accessToken string,
	) (cognitoidentityprovider.GetUserOutput, error)
}

type VerifyTokenResult[PayloadType any] struct {
	Header  CognitoTokenJWTHeader
	Payload PayloadType

	//
	IsValidSignature bool
	IsValidIss       bool
	IsValidExp       bool
	IsValidClientID  bool
	IsValidTokenUse  bool
	Err              error
}

func (result VerifyTokenResult[PayloadType]) IsErr() bool {
	return result.Err != nil
}

func (result VerifyTokenResult[PayloadType]) IsValid() bool {
	return result.IsValidIss &&
		result.IsValidExp &&
		result.IsValidSignature &&
		result.IsValidClientID &&
		result.IsValidTokenUse
}

func NewCognitoClient(ctx context.Context) (CognitoClient, error) {
	// load aws config from environment variables
	cnf, err := NewAWSConfig(ctx)
	if err != nil {
		return nil, errors.LiftWithCtx(err, ctx)
	}
	return &cognitoClient{
		sdkClient: cognitoidentityprovider.NewFromConfig(cnf),
	}, nil
}

type cognitoClient struct {
	sdkClient *cognitoidentityprovider.Client
}

func (client cognitoClient) VerifyAccessToken(
	ctx context.Context,
	userPoolID string,
	token CognitoAccessTokenJWT,
	requestAt time.Time,
) VerifyTokenResult[CognitoAccessTokenJWTPayload] {
	result := VerifyTokenResult[CognitoAccessTokenJWTPayload]{}
	header, payload, err := token.GetHeaderAndPayload()
	if err != nil {
		result.Err = errors.Lift(err)
		return result
	}
	result.Header = header
	result.Payload = payload

	// verify signature
	result.IsValidSignature, err = client.VerifySignature(ctx, token, userPoolID)
	if err != nil {
		result.Err = errors.Lift(err)
		return result
	}
	if !result.IsValidSignature {
		return result
	}

	// expiration
	result.IsValidExp = payload.IsValidExp(requestAt.Unix())

	// iss
	opts := client.sdkClient.Options()
	result.IsValidIss = payload.IsValidIss(
		opts.Region, userPoolID,
	)
	// token use
	result.IsValidTokenUse = payload.IsAccessToken()

	// aud
	result.IsValidClientID = payload.IsValidClientID(
		GetCognitoClientID(),
	)

	// checking kid is slow due to http connection
	// if already invalid, return early
	if !result.IsValid() {
		return result
	}
	return result
}

func (client cognitoClient) GetJWK(userPoolID string) (JSONWebKeys, error) {
	resp, err := http.Get(
		fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json)",
			client.sdkClient.Options().Region, userPoolID,
		),
	)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return JSONWebKeys{}, errors.Lift(err)
	}
	if resp == nil {
		return JSONWebKeys{}, errors.New("http response is empty")
	}
	var keys JSONWebKeys
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return keys, errors.Lift(err)
	}
	return keys, nil
}

func (client cognitoClient) VerifySignature(
	ctx context.Context,
	token CognitoTokener,
	userPoolID string,
) (bool, error) {
	keys, err := client.GetJWK(userPoolID)
	if err != nil {
		return false, errors.Lift(err)
	}
	cognitoToken := token.CognitoTokenJWT()
	header, err := cognitoToken.GetHeader()
	if err != nil {
		return false, errors.Lift(err)
	}
	if header.Algorithm != "RS256" {
		return false, errors.New("unsupported alg, expected RS256")
	}
	jwk, ok := keys.GetByKid(header.KeyID)
	if !ok {
		err = errors.New("failed to get jwk by kid")
		return false, errors.AddTagString(err, "kid", header.KeyID)
	}

	pub, err := jwk.toRSAPublicKey()
	if err != nil {
		return false, errors.Lift(err)
	}
	headerStr, payloadStr, sig := cognitoToken.Parts()
	hash := sha256.Sum256([]byte(headerStr + "." + payloadStr))
	if err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], []byte(sig)); err != nil {
		return false, errors.Lift(err)
	}
	return true, nil
}

func (client cognitoClient) GetUserInfo(
	ctx context.Context,
	accessToken string,
) (cognitoidentityprovider.GetUserOutput, error) {
	output := cognitoidentityprovider.GetUserOutput{}
	rawOutput, err := client.sdkClient.GetUser(ctx, &cognitoidentityprovider.GetUserInput{
		AccessToken: &accessToken,
	})
	if rawOutput != nil {
		output = *rawOutput
	}
	if err != nil {
		return output, errors.Lift(err)
	}
	return output, nil
}
