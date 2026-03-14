package lambdaapigw

import (
	"app/pkg/crosscutting"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponseBody struct {
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

func (body ErrorResponseBody) JsonString() string {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Sprintf(`{"message":"%s"}`, body.Message)
	}
	return string(b)
}

func NewValidateErrorResponseBody(details []ErrorDetail) ErrorResponseBody {
	return ErrorResponseBody{
		Message: "validation error",
		Details: details,
	}
}

func NewValidateErrorResponse(details []ErrorDetail) events.APIGatewayProxyResponse {
	body := NewValidateErrorResponseBody(details)
	return events.APIGatewayProxyResponse{
		StatusCode: crosscutting.APIStatusBadRequest.Int(),
		Body:       body.JsonString(),
	}
}

func NewValidateErrorResponseFromValidateDetails(details []crosscutting.ValidateDetail) events.APIGatewayProxyResponse {
	errorDetails := make([]ErrorDetail, len(details))
	for i, detail := range details {
		errorDetails[i] = ConvertValidateErrorToErrorDetail(detail)
	}
	return NewValidateErrorResponse(errorDetails)
}

func NewBodyDecodingErrorResponseBody(details []ErrorDetail) ErrorResponseBody {
	return ErrorResponseBody{
		Message: "request body decoding error",
		Details: details,
	}
}

func NewBodyDecodingErrorResponse() events.APIGatewayProxyResponse {
	body := NewBodyDecodingErrorResponseBody(nil)
	return events.APIGatewayProxyResponse{
		StatusCode: crosscutting.APIStatusBadRequest.Int(),
		Body:       body.JsonString(),
	}
}

func NewInternalServerErrorResponseBody(details []ErrorDetail) ErrorResponseBody {
	return ErrorResponseBody{
		Message: "internal server error",
		Details: details,
	}
}

func NewInternalServerErrorResponse() events.APIGatewayProxyResponse {
	body := NewInternalServerErrorResponseBody(nil)
	return events.APIGatewayProxyResponse{
		StatusCode: crosscutting.APIStatusInternalServerError.Int(),
		Body:       body.JsonString(),
	}
}

type ErrorDetailType string

type ErrorDetail struct {
	Type    ErrorDetailType `json:"type,omitempty"`
	Field   string          `json:"field,omitempty"`
	Message string          `json:"message"`

	// information
	Min    *int    `json:"min,omitempty"`
	Max    *int    `json:"max,omitempty"`
	Format *string `json:"format,omitempty"`
}

func ConvertInValidTypeToErrorDetailType(invalidType crosscutting.InValidType) ErrorDetailType {
	switch invalidType {
	case crosscutting.InValidTypeRequired:
		return "required"
	case crosscutting.InValidTypeMin:
		return "min"
	case crosscutting.InValidTypeMax:
		return "max"
	case crosscutting.InValidTypeFormat:
		return "format"
	default:
		return "unknown"
	}
}

func ConvertValidateErrorToErrorDetail(detail crosscutting.ValidateDetail) ErrorDetail {
	return ErrorDetail{
		Type:    ConvertInValidTypeToErrorDetailType(detail.Type),
		Field:   detail.Field,
		Message: detail.Message,
		Min:     detail.Min,
		Max:     detail.Max,
		Format:  detail.Format,
	}
}
