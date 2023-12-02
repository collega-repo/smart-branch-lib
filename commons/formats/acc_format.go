package formats

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func AccFormat(formatPattern string, account string) string {
	for i, char := range formatPattern {
		switch char {
		case '.', '-', '_':
			account = fmt.Sprintf(`%s%c%s`, account[:i], char, account[i:])
		}
	}
	return account
}

func StringToDecimal(value string) decimal.Decimal {
	valueDecimal, _ := decimal.NewFromString(value)
	return valueDecimal
}
