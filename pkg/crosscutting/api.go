package crosscutting

type APIStatusCode int

func (value APIStatusCode) IsSuccess() bool {
	return value >= 200 && value < 300
}

func (value APIStatusCode) Int() int {
	return int(value)
}

const (
	APIStatusOK                  APIStatusCode = 200
	APIStatusAccepted            APIStatusCode = 202
	APIStatusBadRequest          APIStatusCode = 400
	APIStatusUnauthorized        APIStatusCode = 401
	APIStatusForbidden           APIStatusCode = 403
	APIStatusNotFound            APIStatusCode = 404
	APIStatusInternalServerError APIStatusCode = 500
)

type InValidType string

func (value InValidType) String() string {
	return string(value)
}

const (
	InValidTypeRequired         InValidType = "required"
	InValidTypeFormat           InValidType = "format"
	InValidTypeNotOption        InValidType = "not_option"
	InValidTypeMin              InValidType = "min"
	InValidTypeMax              InValidType = "max"
	InValidTypeNotInRange       InValidType = "not_in_range"
	InValidTypeMinStrLen        InValidType = "min_string_length"
	InValidTypeMaxStrLen        InValidType = "max_string_length"
	InValidTypeNotInRangeStrLen InValidType = "not_in_range_string_length"
	InValidTypeMaxArrayLen      InValidType = "max_array_length"

	InValidTypeDuplicated  InValidType = "duplicated"
	InvalidTypeContainChar InValidType = "contain_character"
	InValidTypeEmail       InValidType = "email"
)

type ValidateDetail struct {
	Type    InValidType `json:"type,omitempty"`
	Field   string      `json:"field,omitempty"`
	Message string      `json:"message"`

	// information
	Min    *int    `json:"min,omitempty"`
	Max    *int    `json:"max,omitempty"`
	Format *string `json:"format,omitempty"`
}

func NewRequiredValidateDetail(field string) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeRequired,
		Field: field,
	}
}

func NewMinValidateDetail(field string, min int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeMin,
		Field: field,
		Min:   &min,
	}
}

func NewMaxValidateDetail(field string, val int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeMax,
		Field: field,
		Max:   &val,
	}
}

func NewNotInRangeValidateDetail(field string, min int, max int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeNotInRange,
		Field: field,
		Min:   &min,
		Max:   &max,
	}
}

func NewMinStrLenValidateDetail(field string, min int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeMinStrLen,
		Field: field,
		Min:   &min,
	}
}

func NewEmptyStrValidateDetail(field string) ValidateDetail {
	return NewMinStrLenValidateDetail(field, 1)
}

func NewMaxStrLenValidateDetail(field string, val int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeMaxStrLen,
		Field: field,
		Max:   &val,
	}
}

func NewNotInRangeStrLenValidateDetail(field string, min int, max int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeNotInRangeStrLen,
		Field: field,
		Min:   &min,
		Max:   &max,
	}
}

func NEwMaxArrayLenValidateDetail(field string, val int) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeMaxArrayLen,
		Field: field,
		Max:   &val,
	}
}

func NewNotOptionValidateDetail(field string) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeNotOption,
		Field: field,
	}
}

func NewEmailValidateDetail(field string) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeEmail,
		Field: field,
	}
}

func NewContainerCharacterValidateDetail(field string, chars string) ValidateDetail {
	return ValidateDetail{
		Type:   InvalidTypeContainChar,
		Field:  field,
		Format: &chars,
	}
}

func NewDuplicatedValidateDetail(field string) ValidateDetail {
	return ValidateDetail{
		Type:  InValidTypeDuplicated,
		Field: field,
	}
}
