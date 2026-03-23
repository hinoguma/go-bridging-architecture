package lambdaapigw

import "github.com/aws/aws-lambda-go/events"

const RequestIDKey = "request_id"

type HandlerRequest struct {
	Raw events.APIGatewayProxyRequest

	AuthenticatedBankAccountID *string
}

func (req HandlerRequest) GetHeader(key string) (string, bool) {
	if req.Raw.Headers == nil {
		return "", false
	}
	v, ok := req.Raw.Headers[key]
	return v, ok
}

func (req HandlerRequest) GetAuthorization() (string, bool) {

	val, ok := req.GetHeader("Authorization")
	if !ok {
		return "", false
	}
	return val, true
}

func (req HandlerRequest) SetAuthenticatedBankAccountID(id string) HandlerRequest {
	req.AuthenticatedBankAccountID = &id
	return req
}

func (req HandlerRequest) GetAuthenticatedBankAccountID() string {
	if req.AuthenticatedBankAccountID == nil {
		return ""
	}
	return *req.AuthenticatedBankAccountID
}

func NewHandlerRequest(raw events.APIGatewayProxyRequest) HandlerRequest {
	return HandlerRequest{Raw: raw}
}
