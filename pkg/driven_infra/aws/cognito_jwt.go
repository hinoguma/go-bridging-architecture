package driven_infra_aws

import (
	"app/pkg/crosscutting/errors"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const TokenUseAccess string = "access"
const TokenUseID string = "id"
const JWTPartsDelim string = "."
const JWTPartsLength int = 3

func DecodeBase64URL(data string) ([]byte, error) {
	data = strings.Replace(data, "-", "+", -1) // 62nd char of encoding
	data = strings.Replace(data, "_", "/", -1) // 63rd char of encoding

	switch len(data) % 4 { // Pad with trailing '='s
	case 0: // no padding
	case 2:
		data += "==" // 2 pad chars
	case 3:
		data += "=" // 1 pad char
	}
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return b, nil
}

type JSONWebKeys struct {
	Keys []JSONWebKey `json:"keys"`
}

func (model JSONWebKeys) Have(target string) bool {
	_, ok := model.GetByKid(target)
	return ok
}

func (model JSONWebKeys) GetByKid(kid string) (JSONWebKey, bool) {
	for _, key := range model.Keys {
		if key.Kid == kid {
			return key, true
		}
	}
	return JSONWebKey{}, false
}

// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html
type JSONWebKey struct {
	//Algorithm (alg)
	//The alg header parameter represents the cryptographic algorithm that is used to secure the ID token. User pools use an RS256 cryptographic algorithm, which is an RSA signature with SHA-256. For more information on RSA, see RSA cryptography.
	Alg string `json:"alg"`

	// Key ID (kid)
	// The kid is a hint that indicates which key was used to secure the JSON Web Signature (JWS) of the token.
	Kid string `json:"kid"`

	// Key type (kty)
	// The kty parameter identifies the cryptographic algorithm family that is used with the key, such as "RSA" in this example.
	Kty string `json:"kty"`

	//RSA exponent (e)
	//The e parameter contains the exponent value for the RSA public key. It is represented as a Base64urlUInt-encoded value.
	E string `json:"e"`

	//RSA modulus (n)
	//The n parameter contains the modulus value for the RSA public key. It is represented as a Base64urlUInt-encoded value.
	N string `json:"n"`

	//Use (use)
	//The use parameter describes the intended use of the public key. For this example, the use value sig represents signature.
	Use string `json:"use"`
}

func (key JSONWebKey) toRSAPublicKey() (*rsa.PublicKey, error) {
	// key.N, key.E are base64url encoded
	nBytes, err := DecodeBase64URL(key.N)
	if err != nil {
		return nil, errors.Lift(err)
	}
	eBytes, err := DecodeBase64URL(key.E)
	if err != nil {
		return nil, errors.Lift(err)
	}

	n := new(big.Int).SetBytes(nBytes)

	// eBytes is big-endian; convert to int
	e, err := exponentBytesToInt(eBytes)
	if err != nil {
		return nil, errors.Lift(err)
	}
	return &rsa.PublicKey{N: n, E: e}, nil
}

// eBytes を安全に int に変換する（空チェック・正の値・int サイズチェック）
func exponentBytesToInt(eBytes []byte) (int, error) {
	if len(eBytes) == 0 {
		return 0, errors.New("invalid exponent in JWK: empty")
	}
	bi := new(big.Int).SetBytes(eBytes)
	if bi.Sign() <= 0 {
		return 0, errors.New("invalid exponent in JWK: non-positive")
	}
	// 符号ビットを考慮して int に収まるかチェック
	if bi.BitLen() > strconv.IntSize-1 {
		return 0, errors.New("exponent too large for int")
	}
	v := bi.Int64() // bitlen <= IntSize-1 なら Int64 で安全
	if strconv.IntSize == 32 {
		if v > int64(int32(^uint32(0)>>1)) { // max int32
			return 0, errors.New("exponent too large for 32-bit int")
		}
	}
	return int(v), nil
}

/***********************
	Common Token
 ***********************/

type CognitoTokener interface {
	CognitoTokenJWT() CognitoTokenJWT
}

type CognitoTokenJWTHeader struct {
	KeyID     string `json:"kid"`
	Algorithm string `json:"alg"`
}

func NewCognitoTokenJWTHeaderByStr(s string) (CognitoTokenJWTHeader, error) {
	value := CognitoTokenJWTHeader{}
	headerDec, err := DecodeBase64URL(s)
	if err != nil {
		return value, errors.Lift(err)
	}
	err = json.Unmarshal(headerDec, &value)
	if err != nil {
		return value, errors.Lift(err)
	}
	return value, nil
}

type CognitoTokenJWT string

func (value CognitoTokenJWT) Parts() (header, payload, signature string) {
	jwtStr := string(value)
	jwtSlice := strings.Split(jwtStr, JWTPartsDelim)
	if len(jwtSlice) != JWTPartsLength {
		return "", "", ""
	}
	return jwtSlice[0], jwtSlice[1], jwtSlice[2]
}

func (value CognitoTokenJWT) GetHeader() (CognitoTokenJWTHeader, error) {
	header, _, _ := value.Parts()
	return NewCognitoTokenJWTHeaderByStr(header)
}

func (value CognitoTokenJWT) GetSignature() string {
	jwtStr := string(value)
	jwtSlice := strings.Split(jwtStr, JWTPartsDelim)
	if len(jwtSlice) != JWTPartsLength {
		return ""
	}
	return jwtSlice[2]
}

/*
**********************

		ID Token
	 **********************
*/
type CognitoIDTokenJWT CognitoTokenJWT

func (value CognitoIDTokenJWT) CognitoTokenJWT() CognitoTokenJWT {
	return CognitoTokenJWT(value)
}

func (value CognitoIDTokenJWT) GetHeaderAndPayload() (CognitoTokenJWTHeader, CognitoIDTokenJWTPayload, error) {
	headerStr, payloadStr, _ := CognitoTokenJWT(value).Parts()

	header, err := NewCognitoTokenJWTHeaderByStr(headerStr)
	if err != nil {
		return header, CognitoIDTokenJWTPayload{}, errors.Lift(err)
	}

	payload, err := NewCognitoIDTokenJWTPayloadByStr(payloadStr)
	if err != nil {
		return header, payload, errors.Lift(err)
	}
	return header, payload, nil
}

type CognitoIDTokenJWTPayload struct {
	Sub                 string `json:"sub"`
	Aud                 string `json:"aud"`
	EmailVerified       bool   `json:"email_verified"`
	TokenUse            string `json:"token_use"`
	AuthTime            int64  `json:"auth_time"`
	Iss                 string `json:"iss"`
	AuthServiceUserName string `json:"cognito:username"`
	Exp                 int64  `json:"exp"`
	GivenName           string `json:"given_name"`
	Iat                 int64  `json:"iat"`
	Email               string `json:"email"`
}

func (value CognitoIDTokenJWTPayload) IsValidExp(ts int64) bool {
	if value.Exp == 0 {
		return false
	}
	return ts <= value.Exp
}

func (value CognitoIDTokenJWTPayload) IsValidIss(region, poolID string) bool {
	return value.Iss == fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, poolID)
}

func (value CognitoIDTokenJWTPayload) IsIDToken() bool {
	return value.TokenUse == TokenUseID
}

func (value CognitoIDTokenJWTPayload) IsValidAud(cognitoClientID string) bool {
	return value.Aud == cognitoClientID
}

func NewCognitoIDTokenJWTPayloadByStr(s string) (CognitoIDTokenJWTPayload, error) {
	value := CognitoIDTokenJWTPayload{}
	headerDec, err := DecodeBase64URL(s)
	if err != nil {
		return value, errors.Lift(err)
	}
	err = json.Unmarshal(headerDec, &value)
	if err != nil {
		return value, errors.Lift(err)
	}
	return value, nil
}

/***********************
	Access Token
 ***********************/

type CognitoAccessTokenJWT CognitoTokenJWT

func (value CognitoAccessTokenJWT) CognitoTokenJWT() CognitoTokenJWT {
	return CognitoTokenJWT(value)
}

func (value CognitoAccessTokenJWT) GetHeaderAndPayload() (CognitoTokenJWTHeader, CognitoAccessTokenJWTPayload, error) {
	headerStr, payloadStr, _ := CognitoTokenJWT(value).Parts()

	header, err := NewCognitoTokenJWTHeaderByStr(headerStr)
	if err != nil {
		return header, CognitoAccessTokenJWTPayload{}, errors.Lift(err)
	}

	payload, err := NewCognitoAccessTokenJWTPayloadByStr(payloadStr)
	if err != nil {
		return header, payload, errors.Lift(err)
	}
	return header, payload, nil
}

// https://docs.aws.amazon.com/ja_jp/cognito/latest/developerguide/amazon-cognito-user-pools-using-the-access-token.html
/** Sample
{
	"sub":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	"device_key": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	"cognito:groups":[
		"testgroup"
	],
	"iss":"https://cognito-idp.us-west-2.amazonaws.com/us-west-2_example",
	"version":2,
	"client_id":"xxxxxxxxxxxxexample",
	"aud": "https://api.example.com",
	"origin_jti":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	"event_id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	"token_use":"access",
	"scope":"phone openid profile resourceserver.1/appclient2 email",
	"auth_time":1676313851,
	"exp":1676317451,
	"iat":1676313851,
	"jti":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	"username":"my-test-user"
}
*/
type CognitoAccessTokenJWTPayload struct {
	Sub           string   `json:"sub"`
	DeviseKey     string   ` json:"device_key"`
	CognitoGroups []string ` json:"cognito:groups"`
	Iss           string   `json:"iss"`
	Version       int64    ` json:"version"`
	ClientID      string   `json:"client_id"`
	Aud           string   `json:"aud"`
	OriginJti     string   `json:"origin_jti"`
	EventID       string   `json:"event_id"`
	TokenUse      string   `json:"token_use"`
	Scope         string   `json:"scope"`
	AuthTime      int64    `json:"auth_time"`
	Exp           int64    `json:"exp"`
	Iat           int64    `json:"iat"`
	Jti           string   `json:"jti"`
	Username      string   `json:"Username"`
}

func (value CognitoAccessTokenJWTPayload) IsValidExp(ts int64) bool {
	if value.Exp == 0 {
		return false
	}
	return ts <= value.Exp
}

func (value CognitoAccessTokenJWTPayload) IsValidIss(region, poolID string) bool {
	return value.Iss == fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, poolID)
}

func (value CognitoAccessTokenJWTPayload) IsAccessToken() bool {
	return value.TokenUse == TokenUseAccess
}

func (value CognitoAccessTokenJWTPayload) IsValidClientID(cognitoClientID string) bool {
	return value.ClientID == cognitoClientID
}

func NewCognitoAccessTokenJWTPayloadByStr(s string) (CognitoAccessTokenJWTPayload, error) {
	value := CognitoAccessTokenJWTPayload{}
	headerDec, err := DecodeBase64URL(s)
	if err != nil {
		return value, errors.Lift(err)
	}
	err = json.Unmarshal(headerDec, &value)
	if err != nil {
		return value, errors.Lift(err)
	}
	return value, nil
}
