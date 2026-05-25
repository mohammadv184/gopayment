package gopayment

import (
	"fmt"
	"math"
)

// CurrencyCode represents an ISO 4217 currency code or a local one like IRT
type CurrencyCode string

func (cc CurrencyCode) String() string { return string(cc) }

const (
	IRR CurrencyCode = "IRR" // Iranian Rial
	IRT CurrencyCode = "IRT" // Iranian Toman
	USD CurrencyCode = "USD" // US Dollar
	EUR CurrencyCode = "EUR" // Euro
)

// Currency represents a currency with its code and decimal places
type Currency struct {
	Code     CurrencyCode
	Exponent int
}

var currencies = map[CurrencyCode]Currency{
	IRR: {Code: IRR, Exponent: 0},
	IRT: {Code: IRT, Exponent: 0},
	USD: {Code: USD, Exponent: 2},
	EUR: {Code: EUR, Exponent: 2},
}

// GetCurrency returns the Currency for a given code
func GetCurrency(code CurrencyCode) (Currency, error) {
	if c, ok := currencies[code]; ok {
		return c, nil
	}
	return Currency{}, fmt.Errorf("unsupported currency: %s", code)
}

// Amount represents a monetary value in the smallest unit of a currency
type Amount struct {
	Val  int64
	Curr Currency
}

// NewAmount creates a new Amount from units
func NewAmount(value int64, curr Currency) Amount {
	return Amount{Val: value, Curr: curr}
}

// AmountFromFloat creates an Amount from a float64 based on exponent
func AmountFromFloat(val float64, curr Currency) Amount {
	units := int64(math.Round(val * math.Pow10(curr.Exponent)))
	return Amount{Val: units, Curr: curr}
}

func (a Amount) Value() int64       { return a.Val }
func (a Amount) Currency() Currency { return a.Curr }

// ToFloat converts the Amount back to a float64
func (a Amount) ToFloat() float64 {
	return float64(a.Val) / math.Pow10(a.Curr.Exponent)
}

// String returns a formatted string
func (a Amount) String() string {
	return fmt.Sprintf("%.*f %s", a.Curr.Exponent, a.ToFloat(), a.Curr.Code)
}

// ConvertTo handles simple conversions.
// Currently, only IRR -> IRT and IRT -> IRR are supported.
func (a Amount) ConvertTo(target Currency) (Amount, error) {
	if a.Curr.Code == target.Code {
		return a, nil
	}
	if a.Curr.Code == IRT && target.Code == IRR {
		return NewAmount(a.Val*10, target), nil
	}
	if a.Curr.Code == IRR && target.Code == IRT {
		return NewAmount(a.Val/10, target), nil
	}
	return Amount{}, fmt.Errorf("conversion from %s to %s not supported", a.Curr.Code, target.Code)
}
