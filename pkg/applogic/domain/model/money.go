package model

type Currency string

func (value Currency) String() string {
	return string(value)
}

func (value Currency) IsValid() bool {
	switch value {
	case USD, EUR, GBP, JPY:
		return true
	default:
		return false
	}
}

func ValidateCurrency(value string) bool {
	currency := Currency(value)
	return currency.IsValid()
}

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type Money struct {
	Amount   int64
	Currency Currency
}

func (value Money) GreaterThan(other Money) bool {
	if value.Currency != other.Currency {
		return false
	}
	return value.Amount > other.Amount
}
