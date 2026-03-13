package model

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type Money struct {
	Amount   int64
	Currency string
}
