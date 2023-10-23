package formats

import "fmt"

func AccFormat(formatPattern string, account string) string {
	for i, char := range formatPattern {
		switch char {
		case '.', '-', '_':
			account = fmt.Sprintf(`%s%c%s`, account[:i], char, account[i:])
		}
	}
	return account
}
